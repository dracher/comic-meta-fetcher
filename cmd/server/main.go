package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/dracher/comic-meta-fetcher/internal/envutil"
	"github.com/dracher/comic-meta-fetcher/provider/moe"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	envutil.LoadEnv(".env")

	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	p := moe.New()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/search", handleSearch(p))
	mux.HandleFunc("GET /api/comic/{id}", handleFetch(p))

	log.Printf("Server listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, corsMiddleware(mux)))
}

func handleSearch(p *moe.MoeProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			writeError(w, http.StatusBadRequest, "missing required query parameter: q")
			return
		}

		resp, err := p.Search(q)
		if err != nil {
			log.Printf("search error for q=%q: %v", q, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

func handleFetch(p *moe.MoeProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "missing comic id")
			return
		}

		meta, err := p.Fetch(id)
		if err != nil {
			log.Printf("fetch error for id=%q: %v", id, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, meta)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
