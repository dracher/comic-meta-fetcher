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


<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:7510c1e2 -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files

**Architecture in one line:** issues live in a local Dolt DB; sync uses `refs/dolt/data` on your git remote; `.beads/issues.jsonl` is a passive export. See https://github.com/gastownhall/beads/blob/main/docs/SYNC_CONCEPTS.md for details and anti-patterns.

## Session Completion

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
<!-- END BEADS INTEGRATION -->
