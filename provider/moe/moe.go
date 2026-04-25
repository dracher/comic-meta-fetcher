package moe

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracher/comic-meta-fetcher/provider"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://kzz.moe/c/%s.htm"

type MoeProvider struct{}

func New() *MoeProvider {
	return &MoeProvider{}
}

func (p *MoeProvider) Name() string {
	return "moe"
}

func (p *MoeProvider) Fetch(id string) (*provider.ComicMeta, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	url := fmt.Sprintf(baseURL, id)
	doc, err := fetchDocument(url)
	if err != nil {
		return nil, err
	}

	return extractMeta(doc, url, id), nil
}

func MetaFetch(id string) ([]byte, error) {
	p := New()
	meta, err := p.Fetch(id)
	if err != nil {
		return nil, err
	}
	return json.Marshal(meta)
}

func fetchDocument(url string) (*goquery.Document, error) {
	client := provider.NewHTTPClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", provider.DefaultUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}
