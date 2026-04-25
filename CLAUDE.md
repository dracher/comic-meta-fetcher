# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build ./...          # compile all packages
go run ./cmd/comic-meta-fetcher  # fetch and print metadata for the hardcoded test ID
```

No tests or lint config exist yet.

## Architecture

This is a Go CLI tool that fetches comic metadata from multiple sources via a Provider interface pattern.

**`provider/` package** — defines the extensibility contract:
- `Provider` interface: `Name() string` + `Fetch(id string) (*ComicMeta, error)`
- `ComicMeta` — source-agnostic metadata struct (title, authors, tags, rating, etc.). Source-specific fields go into `Extra map[string]any`.
- Shared HTTP helpers (`NewHTTPClient`, `DefaultUserAgent`)

**`provider/<source>/`** — each source is an isolated sub-package implementing `Provider`. Currently only `provider/moe/` (kzz.moe HTML scraper) exists. Adding a new source (e.g. MangaDex) means creating a new sub-package with a struct that satisfies `Provider`.

**`main.go`** — wires a `provider.Provider` and calls `Fetch`, prints JSON. To switch sources, change which provider is instantiated.

## Key conventions

- All providers live under `provider/` as sub-packages named after the source.
- `ComicMeta` common fields use structured types (`[]Author`, `[]Tag`, `*Rating`). Anything source-specific lives in `Extra`.
- The moe provider scrapes HTML via goquery; future providers (e.g. MangaDex) would use JSON APIs instead.
