package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dracher/comic-meta-fetcher/provider"
	"github.com/dracher/comic-meta-fetcher/provider/moe"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <id>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  id  comic ID to fetch (e.g. 50010)")
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	id := flag.Arg(0)

	var p provider.Provider = moe.New()
	fmt.Fprintf(os.Stderr, "Fetching metadata from %s for ID: %s\n", p.Name(), id)

	meta, err := p.Fetch(id)
	if err != nil {
		log.Fatalf("Failed to fetch metadata: %v", err)
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal metadata: %v", err)
	}
	fmt.Println(string(data))
}
