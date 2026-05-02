package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dracher/comic-meta-fetcher/internal/envutil"
	"github.com/dracher/comic-meta-fetcher/provider/moe"
)

func main() {
	envutil.LoadEnv(".env")

	searchQuery := flag.String("s", "", "search for comics by name")
	cookie := flag.String("cookie", "", "cookie header for authenticated requests")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-s query] [-cookie value] [id]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  id        comic ID to fetch (e.g. 50010)")
		fmt.Fprintln(os.Stderr, "  -s        search for comics by name")
		fmt.Fprintln(os.Stderr, "  -cookie   cookie header for authentication")
	}
	flag.Parse()

	var opts []moe.Option
	if *cookie != "" {
		opts = append(opts, moe.WithCookie(*cookie))
	}
	p := moe.New(opts...)

	if *searchQuery != "" {
		fmt.Fprintf(os.Stderr, "Searching %s for: %s\n", p.Name(), *searchQuery)
		resp, err := p.Search(*searchQuery)
		if err != nil {
			log.Fatalf("Search failed: %v", err)
		}
		printJSON(resp)
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	id := flag.Arg(0)

	fmt.Fprintf(os.Stderr, "Fetching metadata from %s for ID: %s\n", p.Name(), id)
	meta, err := p.Fetch(id)
	if err != nil {
		log.Fatalf("Failed to fetch metadata: %v", err)
	}
	printJSON(meta)
}

func printJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}
	fmt.Println(string(data))
}
