# REASONIX.md — comic-meta-fetcher

## Stack

- **Language** — Go 1.26.1 (module `github.com/dracher/comic-meta-fetcher`)
- **HTML scraping** — PuerkitoBio/goquery v1.12.0 (used by the `moe` provider)
- **HTTP** — `net/http` with a shared 30s-timeout client in `provider.NewHTTPClient()`

## Layout

- `cmd/comic-meta-fetcher/main.go` — CLI entry point; wires a Provider, calls Fetch, prints JSON to stdout, logs to stderr
- `provider/provider.go` — `Provider` interface (`Name()`, `Fetch(id)`) + `ComicMeta` struct + shared HTTP helpers
- `provider/moe/` — kzz.moe HTML scraper provider; `moe.go` (HTTP fetch, document loading) and `parser.go` (HTML → ComicMeta extraction via goquery + regexp)

## Commands

```bash
go build ./...                          # compile all packages
go run ./cmd/comic-meta-fetcher <id>    # fetch and print JSON for a comic ID
```

## Conventions

- Providers live under `provider/<source>/` as standalone packages implementing the `Provider` interface.
- Structured metadata goes into `ComicMeta` typed fields; source-specific data goes in `Extra map[string]any`.
- JSON output uses `json.MarshalIndent` with 2-space indent. Errors and progress messages go to stderr.

## Watch out for

- **No tests exist** — any edit to the provider or parser has no CI gate catching regressions.
- **The `moe` provider is a fragile HTML scraper** — it depends on specific CSS class names (`text_bglight_big`, `img_book`, `book_score`) and regex that extracts JS variable assignments from inline `<script>` blocks. If kzz.moe changes its markup, the parser silently returns empty fields.
- **Only one provider is wired in `main.go`** — to switch sources you must edit the file directly; there is no provider registry or config-driven selection.
