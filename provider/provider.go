package provider

import (
	"net/http"
	"time"
)

type ComicMeta struct {
	Source      string         `json:"source"`
	ID          string         `json:"id"`
	URL         string         `json:"url,omitempty"`
	Title       string         `json:"title"`
	Aliases     []string       `json:"aliases,omitempty"`
	Description string         `json:"description,omitempty"`
	CoverImage  string         `json:"cover_image,omitempty"`
	Authors     []Author       `json:"authors,omitempty"`
	Status      string         `json:"status,omitempty"`
	Language    string         `json:"language,omitempty"`
	Genres      []string       `json:"genres,omitempty"`
	Tags        []Tag          `json:"tags,omitempty"`
	Rating      *Rating        `json:"rating,omitempty"`
	FollowCount int            `json:"follow_count,omitempty"`
	Extra       map[string]any `json:"extra,omitempty"`
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

type Provider interface {
	Name() string
	Fetch(id string) (*ComicMeta, error)
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
