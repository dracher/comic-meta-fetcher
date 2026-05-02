package moe

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/dracher/comic-meta-fetcher/provider"

	"github.com/PuerkitoBio/goquery"
)

var (
	// Match actual disp_divinfo calls (start with "div_info_"+"N"), skipping the function definition.
	reDispDivInfo = regexp.MustCompile(`(?s)disp_divinfo\s*\(\s*"div_info_"\s*\+\s*"\d+"\s*,(.+?)\)\s*;`)
	reDispDivPage = regexp.MustCompile(`(?s)disp_divpage\s*\(\s*"[^"]*"\s*,\s*"[^"]*"\s*,\s*"(\d+)"`)
	reQuotedStr   = regexp.MustCompile(`"([^"]*)"`)
	reHTMLTag     = regexp.MustCompile(`<[^>]+>`)
)

func parseSearchResults(doc *goquery.Document, query string) (*provider.SearchResponse, error) {
	resp := &provider.SearchResponse{
		Query: query,
		Page:  1,
	}

	var scriptContent strings.Builder
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptContent.WriteString(s.Text())
	})
	content := scriptContent.String()

	if m := reDispDivPage.FindStringSubmatch(content); len(m) >= 2 {
		resp.TotalPages, _ = strconv.Atoi(m[1])
	}

	// Each match group 1 contains the 12 params after div_id:
	// book_url, cover_url, border_color,
	// tag_jp, tag_en, tag_end, tag_brk,
	// score, name, author, volume_status, update_date
	for _, match := range reDispDivInfo.FindAllStringSubmatch(content, -1) {
		params := reQuotedStr.FindAllStringSubmatch(match[1], -1)
		if len(params) < 12 {
			continue
		}

		bookURL := params[0][1]
		coverURL := params[1][1]
		tagJP := params[3][1]
		tagEN := params[4][1]
		tagEnd := params[5][1]
		tagBrk := params[6][1]
		score := params[7][1]
		name := reHTMLTag.ReplaceAllString(params[8][1], "")
		author := params[9][1]
		volumeStatus := params[10][1]
		updateDate := params[11][1]

		scoreFloat, _ := strconv.ParseFloat(score, 64)

		var tags []string
		if tagJP == "" {
			tags = append(tags, "日語")
		}
		if tagEN == "" {
			tags = append(tags, "英文")
		}
		if tagEnd == "" {
			tags = append(tags, "完結")
		}
		if tagBrk == "" {
			tags = append(tags, "停更")
		}

		resp.Results = append(resp.Results, provider.SearchResult{
			Source:     "moe",
			ID:        extractIDFromURL(bookURL),
			URL:       bookURL,
			Title:     name,
			Authors:   author,
			CoverImage: coverURL,
			Score:     scoreFloat,
			Status:    volumeStatus,
			UpdateDate: updateDate,
			Tags:      tags,
		})
	}

	return resp, nil
}

func extractIDFromURL(u string) string {
	idx := strings.LastIndex(u, "/")
	if idx == -1 {
		return u
	}
	return strings.TrimSuffix(u[idx+1:], ".htm")
}
