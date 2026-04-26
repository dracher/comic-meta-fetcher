# Comic Meta Fetcher

A Go CLI tool that fetches comic/book metadata from multiple online sources through a unified provider interface.

## Features

- Fetch comic metadata such as title, authors, tags, rating, cover image, and more
- Pluggable architecture — easily add new data sources by implementing the `Provider` interface
- Structured output in JSON format
- Currently supports [kzz.moe](https://kzz.moe) as a data source

## Usage

```bash
# Build the project
go build ./...

# Fetch metadata for a specific comic ID
go run ./cmd/comic-meta-fetcher <id>
```

### Example

```bash
$ go run ./cmd/comic-meta-fetcher 50010
```

Output (pretty-printed JSON):

```json
{
  "source": "moe",
  "id": "50010",
  "url": "https://kzz.moe/c/50010.htm",
  "title": "...",
  "authors": [
    {
      "name": "...",
      "url": "...",
      "role": "author"
    }
  ],
  "tags": [...],
  "rating": {
    "score": 4.5,
    "score_count": 1234,
    "score_distribution": {
      "5星": "65%",
      "4星": "20%"
    }
  },
  ...
}
```

## Installation

```bash
# Clone the repository
git clone https://github.com/dracher/comic-meta-fetcher.git
cd comic-meta-fetcher

# Build the binary
go build -o comic-meta-fetcher ./cmd/comic-meta-fetcher

# Run (requires Go 1.26+)
./comic-meta-fetcher <id>
```

## Architecture

The project follows a provider-based plugin architecture, making it easy to add new data sources.

```
comic-meta-fetcher/
├── cmd/
│   └── comic-meta-fetcher/
│       └── main.go          # CLI entry point, wires provider and outputs JSON
├── provider/
│   ├── provider.go          # Provider interface & shared types
│   └── moe/
│       ├── moe.go           # kzz.moe provider implementation
│       └── parser.go        # HTML scraping & parsing logic
├── CLAUDE.md                # Guidance for Claude Code
├── go.mod
└── go.sum
```

### Core Types

#### `Provider` Interface

```go
type Provider interface {
    Name() string
    Fetch(id string) (*ComicMeta, error)
}
```

Every data source implements this interface. `Name()` returns the source identifier, and `Fetch(id)` returns structured metadata.

#### `ComicMeta` Struct

A source-agnostic metadata container with reusable fields:

| Field         | Type                | Description                     |
|---------------|---------------------|---------------------------------|
| `Source`      | `string`            | Data source name                |
| `ID`          | `string`            | Comic identifier                |
| `URL`         | `string`            | Source page URL                 |
| `Title`       | `string`            | Comic title                     |
| `Aliases`     | `[]string`          | Alternative titles              |
| `Description` | `string`            | Comic description               |
| `CoverImage`  | `string`            | URL of the cover image          |
| `Authors`     | `[]Author`          | List of authors                 |
| `Status`      | `string`            | Publication status              |
| `Language`    | `string`            | Language                        |
| `Genres`      | `[]string`          | Genre list                      |
| `Tags`        | `[]Tag`             | Tag list with counts            |
| `Rating`      | `*Rating`           | Score, count, and distribution  |
| `FollowCount` | `int`               | Number of followers             |
| `Extra`       | `map[string]any`    | Source-specific metadata        |

## Extending: Adding a New Provider

1. Create a new package under `provider/`, e.g. `provider/mangadex/`
2. Define a struct that implements `provider.Provider`
3. Wire it in `main.go` by replacing the instantiated provider

```go
// provider/mangadex/mangadex.go

package mangadex

import "github.com/dracher/comic-meta-fetcher/provider"

type MangadexProvider struct{}

func New() *MangadexProvider {
    return &MangadexProvider{}
}

func (p *MangadexProvider) Name() string {
    return "mangadex"
}

func (p *MangadexProvider) Fetch(id string) (*provider.ComicMeta, error) {
    // Fetch & parse data from MangaDex API
    // Return a populated *provider.ComicMeta
}
```

## Dependencies

- [goquery](https://github.com/PuerkitoBio/goquery) — HTML parsing and selection (used by the `moe` provider)

## License

MIT