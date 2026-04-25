package moe

import (
	"fmt"
	"net/http"

	"github.com/dracher/comic-meta-fetcher/provider"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://kzz.moe/c/%s.htm"

type MoeProvider struct {
	client *http.Client
}

func New() *MoeProvider {
	return &MoeProvider{client: provider.NewHTTPClient()}
}

func (p *MoeProvider) Name() string {
	return "moe"
}

func (p *MoeProvider) Fetch(id string) (*provider.ComicMeta, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	url := fmt.Sprintf(baseURL, id)
	doc, err := p.fetchDocument(url)
	if err != nil {
		return nil, err
	}

	return extractMeta(doc, url, id), nil
}

func (p *MoeProvider) fetchDocument(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", provider.DefaultUserAgent())

	resp, err := p.client.Do(req)
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
