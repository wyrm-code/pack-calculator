package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/example/pack-calculator/internal/solver"
)

type calcRequest struct {
	Items int   `json:"items"`
	Sizes []int `json:"sizes,omitempty"`
}

type packOut struct {
	Size int `json:"size"`
	Qty  int `json:"qty"`
}

type calcResponse struct {
	ItemsOrdered int       `json:"itemsOrdered"`
	TotalItems   int       `json:"totalItems"`
	Packs        []packOut `json:"packs"`
}

func parseSizesFromEnv() ([]int, error) {
	raw := os.Getenv("PACK_SIZES")
	if raw == "" {
		return nil, fmt.Errorf("PACK_SIZES is empty; provide sizes in request or set env var")
	}
	return parseSizesCSV(raw)
}

func parseSizesCSV(csv string) ([]int, error) {
	if strings.TrimSpace(csv) == "" {
		return nil, fmt.Errorf("sizes string empty")
	}
	parts := strings.Split(csv, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil || v <= 0 {
			return nil, fmt.Errorf("invalid size: %q", p)
		}
		out = append(out, v)
	}
	return out, nil
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	var req calcRequest

	switch r.Method {
	case http.MethodGet:
		itemsStr := r.URL.Query().Get("items")
		if itemsStr == "" {
			http.Error(w, "missing items", http.StatusBadRequest)
			return
		}
		items, err := strconv.Atoi(itemsStr)
		if err != nil || items <= 0 {
			http.Error(w, "items must be a positive integer", http.StatusBadRequest)
			return
		}
		req.Items = items

		sizesStr := r.URL.Query().Get("sizes")
		if sizesStr != "" {
			sz, err := parseSizesCSV(sizesStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			req.Sizes = sz
		}

	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if req.Items <= 0 {
		http.Error(w, "items must be > 0", http.StatusBadRequest)
		return
	}
	if len(req.Sizes) == 0 {
		sizes, err := parseSizesFromEnv()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Sizes = sizes
	}

	counts, total, err := solver.Solve(req.Items, req.Sizes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf(
		"calc handled items=%d sizes=%v total=%d packs=%v",
		req.Items, req.Sizes, total, counts,
	)

	// order packs desc for nicer output
	packs := make([]packOut, 0, len(counts))
	for _, s := range solver.SortedDescKeys(counts) {
		packs = append(packs, packOut{Size: s, Qty: counts[s]})
	}

	resp := calcResponse{
		ItemsOrdered: req.Items,
		TotalItems:   total,
		Packs:        packs,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calc", calcHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", uiHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
