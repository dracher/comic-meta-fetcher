package moe

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/dracher/comic-meta-fetcher/provider"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL   = "https://kzz.moe/c/%s.htm"
	searchURL = "https://kzz.moe/list.php"
)

type Option func(*MoeProvider)

func WithCookie(cookie string) Option {
	return func(p *MoeProvider) {
		p.cookie = cookie
	}
}

type MoeProvider struct {
	client *http.Client
	cookie string
}

func New(opts ...Option) *MoeProvider {
	p := &MoeProvider{
		client: provider.NewHTTPClient(),
		cookie: os.Getenv("MOE_COOKIE"),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *MoeProvider) Name() string {
	return "moe"
}

func (p *MoeProvider) Fetch(ctx context.Context, id string) (*provider.ComicMeta, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	u := fmt.Sprintf(baseURL, id)
	doc, err := p.fetchDocument(ctx, u)
	if err != nil {
		return nil, err
	}

	return extractMeta(doc, u, id), nil
}

func (p *MoeProvider) Search(ctx context.Context, query string) (*provider.SearchResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	u := searchURL + "?s=" + url.QueryEscape(query)
	doc, err := p.fetchDocument(ctx, u)
	if err != nil {
		return nil, err
	}

	return parseSearchResults(doc, query)
}

func (p *MoeProvider) fetchDocument(ctx context.Context, u string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", provider.DefaultUserAgent())
	if p.cookie != "" {
		req.Header.Set("Cookie", p.cookie)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", u, err)
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
