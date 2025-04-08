package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/seamusr7/go-webcrawler/crawler"
	"github.com/seamusr7/go-webcrawler/report"
)

// CrawlRequest represents the expected JSON body for crawl API
type CrawlRequest struct {
	URL      string `json:"url"`
	MaxPages int    `json:"maxPages"`
}

// StartServer initializes and starts the REST API server
func StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/api/crawl", handleCrawlPost).Methods("POST")
	router.HandleFunc("/api/crawl", handleCrawlGet).Methods("GET")

	log.Println("🚀 Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleCrawlPost handles POST /api/crawl with JSON body
func handleCrawlPost(w http.ResponseWriter, r *http.Request) {
	var req CrawlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" || req.MaxPages <= 0 {
		http.Error(w, "URL and maxPages must be provided and valid", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Starting crawl on: %s (max %d pages)", req.URL, req.MaxPages)

	pages := crawler.StartCrawling(req.URL, req.MaxPages)
	format := r.URL.Query().Get("format")

	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=seo_report.csv")
		report.ExportToCSV(w, pages)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}

// handleCrawlGet handles GET /api/crawl with query params
func handleCrawlGet(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	maxPagesStr := r.URL.Query().Get("maxPages")
	format := r.URL.Query().Get("format")

	if url == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	maxPages, err := strconv.Atoi(maxPagesStr)
	if err != nil || maxPages <= 0 {
		maxPages = 10 // Default fallback
	}

	log.Printf("🔍 Starting crawl on: %s (max %d pages)", url, maxPages)

	pages := crawler.StartCrawling(url, maxPages)

	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=seo_report.csv")
		report.ExportToCSV(w, pages)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}
