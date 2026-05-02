# Comic Meta Fetcher

A Go tool that fetches comic/book metadata from multiple online sources through a unified provider interface. Provides both a CLI and an HTTP API server.

## Features

- Fetch comic metadata (title, authors, tags, rating, cover image, etc.) by ID
- Search comics by name
- Pluggable architecture — add new data sources by implementing the `Provider` interface
- Structured JSON output
- CLI and HTTP API server
- Currently supports [kzz.moe](https://kzz.moe) as a data source

## Setup

```bash
git clone https://github.com/dracher/comic-meta-fetcher.git
cd comic-meta-fetcher
go build ./...
```

### Configuration

Create a `.env` file in the project root for authentication:

```
MOE_COOKIE="your_cookie_value_here"
```

The cookie is loaded automatically. You can also pass it via `-cookie` flag or `MOE_COOKIE` environment variable.

## CLI Usage

```bash
# Fetch metadata by comic ID
go run ./cmd/cli <id>

# Search comics by name
go run ./cmd/cli -s "七龙珠"

# Override cookie
go run ./cmd/cli -s "七龙珠" -cookie "VOLSESS=..."
```

### Example

```bash
$ go run ./cmd/cli 50010
```

```json
{
  "source": "moe",
  "id": "50010",
  "title": "...",
  "authors": [{ "name": "...", "role": "author" }],
  "rating": { "score": 4.5, "score_count": 1234 }
}
```

```bash
$ go run ./cmd/cli -s "七龙珠"
```

```json
{
  "query": "七龙珠",
  "total_pages": 1,
  "page": 1,
  "results": [
    { "id": "11283", "title": "七龍珠", "authors": "鳥山明", "score": 9.6, "tags": ["完結"] }
  ]
}
```

## HTTP API Server

```bash
go run ./cmd/server              # listen on :8080
go run ./cmd/server -addr :3000  # custom port
```

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/search?q=<query>` | Search comics by name |
| `GET` | `/api/comic/{id}` | Fetch comic metadata by ID |

### Examples

```bash
curl 'http://localhost:8080/api/search?q=七龙珠'
curl 'http://localhost:8080/api/comic/11283'
```

Error responses return `{"error": "..."}` with appropriate HTTP status codes.

## Architecture

```
comic-meta-fetcher/
├── cmd/
│   ├── cli/main.go              # CLI entry point
│   └── server/main.go           # HTTP API server
├── internal/
│   └── envutil/envutil.go       # .env file loader
├── provider/
│   ├── provider.go              # Provider/Searcher interfaces & shared types
│   └── moe/
│       ├── moe.go               # kzz.moe provider (fetch + search)
│       ├── parser.go            # HTML scraping for comic detail pages
│       └── search.go            # JS-embedded search result parser
├── CLAUDE.md
├── go.mod
└── go.sum
```

### Core Interfaces

```go
type Provider interface {
    Name() string
    Fetch(id string) (*ComicMeta, error)
}

type Searcher interface {
    Search(query string) (*SearchResponse, error)
}
```

### Key Types

| Type | Description |
|------|-------------|
| `ComicMeta` | Full metadata: title, authors, tags, rating, description, etc. |
| `SearchResult` | Lightweight result: title, authors, score, status, tags |
| `SearchResponse` | Paginated search results with query and total pages |

## Extending: Adding a New Provider

1. Create a new package under `provider/`, e.g. `provider/mangadex/`
2. Implement `provider.Provider` (and optionally `provider.Searcher`)
3. Wire it in `cmd/cli/main.go` or `cmd/server/main.go`

## Dependencies

- [goquery](https://github.com/PuerkitoBio/goquery) — HTML parsing (used by the `moe` provider)

## License

MIT
