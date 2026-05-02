package provider

import (
	"net/http"
	"time"
)

type ComicMeta struct {
	Source          string         `json:"source"`
	ID              string         `json:"id"`
	URL             string         `json:"url,omitempty"`
	Title           string         `json:"title"`
	Aliases         []string       `json:"aliases,omitempty"`
	Description     string         `json:"description,omitempty"`
	CoverImage      string         `json:"cover_image,omitempty"`
	CoverImageLocal string         `json:"cover_image_local,omitempty"`
	Authors         []Author       `json:"authors,omitempty"`
	Status          string         `json:"status,omitempty"`
	Language        string         `json:"language,omitempty"`
	Genres          []string       `json:"genres,omitempty"`
	Tags            []Tag          `json:"tags,omitempty"`
	Rating          *Rating        `json:"rating,omitempty"`
	FollowCount     int            `json:"follow_count,omitempty"`
	Type            string         `json:"type"`
	InLibrary       bool           `json:"in_library,omitempty"`
	Blocked         bool           `json:"blocked,omitempty"`
	CreatedAt       string         `json:"created_at,omitempty"`
	UpdatedAt       string         `json:"updated_at,omitempty"`
	Extra           map[string]any `json:"extra,omitempty"`
}

type Author struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
	Role string `json:"role,omitempty"`
}

type Tag struct {
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"`
}

type Rating struct {
	Score      float64          `json:"score,omitempty"`
	ScoreCount int              `json:"score_count,omitempty"`
	ScoreDist  map[string]string `json:"score_distribution,omitempty"`
}

type SearchResult struct {
	Source     string   `json:"source"`
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Title     string   `json:"title"`
	Authors   string   `json:"authors"` // raw comma-separated string; search results lack structured author data
	CoverImage string  `json:"cover_image,omitempty"`
	Score     float64  `json:"score"`
	Status    string   `json:"status,omitempty"`
	UpdateDate string  `json:"update_date,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type SearchResponse struct {
	Query      string         `json:"query"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	Results    []SearchResult `json:"results"`
}

type Provider interface {
	Name() string
	Fetch(id string) (*ComicMeta, error)
}

type Searcher interface {
	Search(query string) (*SearchResponse, error)
}

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func DefaultUserAgent() string {
	return defaultUserAgent
}
