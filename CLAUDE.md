# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build ./...                   # compile all packages
go vet ./...                     # static analysis
go run ./cmd/cli <id>            # fetch metadata by ID
go run ./cmd/cli -s "query"      # search by name
go run ./cmd/server              # start HTTP API on :8080
```

No tests or lint config exist yet.

## Architecture

This is a Go tool that fetches comic metadata from multiple sources via a Provider interface pattern. It has two entry points: a CLI (`cmd/cli`) and an HTTP API server (`cmd/server`).

**`internal/envutil/`** — shared `.env` file loader, used by both CLI and server.

**`provider/` package** — defines the extensibility contracts:
- `Provider` interface: `Name() string` + `Fetch(id string) (*ComicMeta, error)`
- `Searcher` interface: `Search(query string) (*SearchResponse, error)`
- `ComicMeta` — source-agnostic metadata struct (title, authors, tags, rating, etc.). Source-specific fields go into `Extra map[string]any`.
- `SearchResult` / `SearchResponse` — lightweight search result types with pagination.
- Shared HTTP helpers (`NewHTTPClient`, `DefaultUserAgent`)

**`provider/<source>/`** — each source is an isolated sub-package implementing `Provider` and `Searcher`. Currently only `provider/moe/` (kzz.moe) exists:
- `moe.go` — provider struct, `Fetch` and `Search` methods, cookie/option support
- `parser.go` — HTML scraper for comic detail pages using goquery
- `search.go` — regex-based parser for JS-embedded search results (`disp_divinfo` calls)

**`cmd/cli/`** — CLI entry point. Supports `-s` for search, positional arg for fetch, `-cookie` for auth override.

**`cmd/server/`** — HTTP API server with `GET /api/search?q=` and `GET /api/comic/{id}`.

## Key conventions

- All providers live under `provider/` as sub-packages named after the source.
- `ComicMeta` common fields use structured types (`[]Author`, `[]Tag`, `*Rating`). Anything source-specific lives in `Extra`.
- The moe provider scrapes HTML via goquery for detail pages; search results are parsed from JS function calls via regex.
- Cookie defaults to `MOE_COOKIE` env var (loaded from `.env`), overridable via `WithCookie` option.
- Shared utilities go in `internal/`.
