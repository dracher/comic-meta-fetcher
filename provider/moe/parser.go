package moe

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/dracher/comic-meta-fetcher/provider"

	"github.com/PuerkitoBio/goquery"
)

var (
	reCategoryEntry = regexp.MustCompile(`([\p{Han}]+)\s*\((\d+)\)`)
	reScoreCount    = regexp.MustCompile(`(\d+)人評價`)
	reScoreStar     = regexp.MustCompile(`([１２３４５])星\s*.*?(\d+%)`)
	reDesc          = regexp.MustCompile(`document\.getElementById\("div_desc_content"\)\.innerHTML\s*=\s*"(.*?)";`)

	reJSVarStr map[string]*regexp.Regexp
	reJSVarInt map[string]*regexp.Regexp
)

func init() {
	reJSVarStr = make(map[string]*regexp.Regexp, len(jsVarNames))
	reJSVarInt = make(map[string]*regexp.Regexp, len(jsVarNames))
	for _, name := range jsVarNames {
		reJSVarStr[name] = regexp.MustCompile(`var\s+` + name + `\s*=\s*"([^"]*)"`)
		reJSVarInt[name] = regexp.MustCompile(`var\s+` + name + `\s*=\s*parseInt\(\s*"([^"]*)"\s*\)`)
	}
}

func extractMeta(doc *goquery.Document, pageURL, bookID string) *provider.ComicMeta {
	meta := &provider.ComicMeta{
		Source: "moe",
		ID:     bookID,
		URL:    pageURL,
		Extra:  make(map[string]any),
	}

	// Title
	meta.Title = strings.TrimSpace(doc.Find("font.text_bglight_big").Text())
	if meta.Title == "" {
		meta.Title, _ = doc.Find("img.img_book").Attr("alt")
	}

	// Cover image
	meta.CoverImage, _ = doc.Find("img.img_book").Attr("src")

	// Author
	authorLink := doc.Find("td.author a[href*='list.php?s=']").First()
	authorName := strings.TrimSpace(authorLink.Text())
	authorURL, _ := authorLink.Attr("href")
	if authorName != "" {
		meta.Authors = append(meta.Authors, provider.Author{
			Name: authorName,
			URL:  authorURL,
			Role: "author",
		})
	}

	// Aliases
	aliasesText := ""
	doc.Find("td.author").Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "font" {
			class, _ := s.Attr("class")
			text := strings.TrimSpace(s.Text())
			if class == "text_bglight" && strings.HasPrefix(text, "(") {
				aliasesText = text
			}
		}
	})
	if aliasesText != "" {
		meta.Aliases = parseAliases(aliasesText)
	}

	// Status info
	doc.Find("td.author font.text_bglight").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "狀態：") {
			meta.Status = extractValue(text, "狀態：", "地區")
			meta.Extra["region"] = extractValue(text, "地區：", "語言")
			meta.Language = extractValue(text, "語言：", "最後出版")
			meta.Extra["last_publish"] = extractValue(text, "最後出版：", "更新")
		}
	})

	// Version info
	doc.Find("td.author font.text_bglight").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "版本：") {
			meta.Extra["version"] = extractValue(text, "版本：", "掃者")
			meta.Extra["scanner"] = extractValue(text, "掃者：", "維護者")
		}
	})

	// Maintainer
	maintainerLink := doc.Find("td.author font.text_bglight a[href*='u/']").First()
	if m := strings.TrimSpace(maintainerLink.Text()); m != "" {
		meta.Extra["maintainer"] = m
		meta.Extra["maintainer_url"], _ = maintainerLink.Attr("href")
	}

	// Statistics
	doc.Find("td.author font.text_bglight").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "訂閱：") {
			meta.FollowCount = extractInt(text, "訂閱：", "收藏")
			meta.Extra["favorites"] = extractInt(text, "收藏：", "讀過")
			meta.Extra["read_count"] = extractInt(text, "讀過：", "熱度")
			meta.Extra["heat"] = extractInt(text, "熱度：", "")
		}
	})

	// Categories/tags
	doc.Find("td.author font.text_bglight").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "分類：") {
			catText := strings.TrimPrefix(text, "分類：")
			for _, match := range reCategoryEntry.FindAllStringSubmatch(catText, -1) {
				if len(match) >= 3 {
					count, _ := strconv.Atoi(match[2])
					name := strings.TrimSpace(match[1])
					meta.Tags = append(meta.Tags, provider.Tag{Name: name, Count: count})
					meta.Genres = append(meta.Genres, name)
				}
			}
		}
	})

	// Score/rating
	scoreText := strings.TrimSpace(doc.Find("table.book_score td font").First().Text())
	if scoreText != "" {
		rating := &provider.Rating{
			ScoreDist: make(map[string]string),
		}
		rating.Score, _ = strconv.ParseFloat(scoreText, 64)

		scoreCountText := doc.Find("table.book_score td font.font_size_s").Text()
		if m := reScoreCount.FindStringSubmatch(scoreCountText); len(m) >= 2 {
			rating.ScoreCount, _ = strconv.Atoi(m[1])
		}

		doc.Find("table.book_score font.scorestar").Each(func(i int, s *goquery.Selection) {
			for _, line := range strings.Split(s.Text(), "\n") {
				line = strings.TrimSpace(line)
				if !strings.Contains(line, "星") {
					continue
				}
				m := reScoreStar.FindStringSubmatch(line)
				if len(m) >= 3 {
					rating.ScoreDist[m[1]+"星"] = m[2]
				}
			}
		})
		meta.Rating = rating
	}

	// Description and JS variables from scripts
	jsVars := make(map[string]string)
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		text := s.Text()

		if m := reDesc.FindStringSubmatch(text); len(m) >= 2 {
			desc := m[1]
			desc = strings.ReplaceAll(desc, `<br />`, "\n")
			desc = strings.ReplaceAll(desc, `<br>`, "\n")
			meta.Description = html.UnescapeString(desc)
		}

		for _, varName := range jsVarNames {
			if m := reJSVarStr[varName].FindStringSubmatch(text); len(m) >= 2 {
				jsVars[varName] = m[1]
				continue
			}
			if m := reJSVarInt[varName].FindStringSubmatch(text); len(m) >= 2 {
				jsVars[varName] = m[1]
			}
		}
	})
	if len(jsVars) > 0 {
		meta.Extra["js_variables"] = jsVars
	}

	// Meta tags
	metaTags := make(map[string]string)
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		property, _ := s.Attr("property")
		content, _ := s.Attr("content")
		key := name
		if key == "" {
			key = property
		}
		if key != "" && content != "" {
			metaTags[key] = content
		}
	})
	if len(metaTags) > 0 {
		meta.Extra["meta_tags"] = metaTags
	}

	return meta
}

var jsVarNames = []string{
	"bookid", "bookstatus", "is_jpn", "is_eng", "is_hd",
	"is_color", "is_r18", "is_blocked", "is_internal", "is_watermark",
}

func parseAliases(raw string) []string {
	var result []string
	if strings.HasPrefix(raw, "(") {
		if idx := strings.Index(raw, ")"); idx != -1 {
			result = append(result, raw[1:idx])
			raw = raw[idx+1:]
		}
	}
	for _, part := range strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '　'
	}) {
		if p := strings.TrimSpace(part); p != "" {
			result = append(result, p)
		}
	}
	return result
}

func extractValue(text, start, end string) string {
	idx := strings.Index(text, start)
	if idx == -1 {
		return ""
	}
	val := text[idx+len(start):]
	if end != "" {
		if endIdx := strings.Index(val, end); endIdx != -1 {
			val = val[:endIdx]
		}
	}
	return strings.TrimSpace(val)
}

func extractInt(text, start, end string) int {
	val := extractValue(text, start, end)
	val = strings.ReplaceAll(val, ",", "")
	val = strings.ReplaceAll(val, " ", "")
	val = strings.ReplaceAll(val, "　", "")
	num, _ := strconv.Atoi(val)
	return num
}
