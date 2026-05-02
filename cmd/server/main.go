package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dracher/comic-meta-fetcher/provider/moe"
)

func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

var provider *moe.MoeProvider

func main() {
	loadEnv(".env")

	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	provider = moe.New()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/search", handleSearch)
	mux.HandleFunc("GET /api/comic/{id}", handleFetch)

	log.Printf("Server listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "missing required query parameter: q")
		return
	}

	resp, err := provider.Search(q)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing comic id")
		return
	}

	meta, err := provider.Fetch(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, meta)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
	fmt.Fprintln(os.Stderr, "Comic Meta Fetcher API Server")
}
