package main

import (
	"encoding/json"
	"fmt"
	"log"

	"comic-meta-fetcher/provider"
	"comic-meta-fetcher/provider/moe"
)

func main() {
	id := "50010"

	var p provider.Provider = moe.New()
	fmt.Printf("Fetching metadata from %s for ID: %s\n", p.Name(), id)

	meta, err := p.Fetch(id)
	if err != nil {
		log.Fatalf("Failed to fetch metadata: %v", err)
	}

	data, _ := json.MarshalIndent(meta, "", "  ")
	fmt.Println(string(data))
}
