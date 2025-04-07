// main.go - SEO Web Crawler
// Entry point for SEO-focused web crawler using modular structure.

package main

import (
	"fmt"
	"sync"

	"github.com/seamusr7/go-webcrawler/crawler"
	"github.com/seamusr7/go-webcrawler/report"
)

func main() {
	startURL := "https://golang.org"
	maxPages := 50

	jobChan := make(chan [2]string, 100)
	pageResults := make(chan crawler.PageInfo, 100)
	linkResults := make(chan []string, 100)
	done := make(chan struct{})

	var wg sync.WaitGroup
	var managerWg sync.WaitGroup

	var allPages []crawler.PageInfo
	var resultsMu sync.Mutex

	visited := make(map[string]bool)
	var visitedMu sync.Mutex

	// Manager goroutine for handling results and new links
	managerWg.Add(1)
	go func() {
		defer managerWg.Done()
		for {
			select {
			case links, ok := <-linkResults:
				if !ok {
					return
				}
				for _, link := range links {
					visitedMu.Lock()
					if !visited[link] && len(allPages) < maxPages {
						visited[link] = true
						wg.Add(1)
						jobChan <- [2]string{link, "referrer unknown"}
					}
					visitedMu.Unlock()
				}
			case page, ok := <-pageResults:
				if !ok {
					return
				}
				resultsMu.Lock()
				if len(allPages) < maxPages {
					allPages = append(allPages, page)
				}
				resultsMu.Unlock()
			case <-done:
				return
			}
		}
	}()

	// Start crawling from the initial URL
	visitedMu.Lock()
	visited[startURL] = true
	visitedMu.Unlock()
	wg.Add(1)
	jobChan <- [2]string{startURL, ""}

	// Launch worker goroutines
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go crawler.Worker(i, jobChan, pageResults, linkResults, &wg)
	}

	// Wait for all crawl jobs to complete
	wg.Wait()
	close(done)
	managerWg.Wait()

	close(jobChan)
	close(linkResults)
	close(pageResults)

	fmt.Println("\n✅ Done crawling.")
	report.Generate(allPages)
}
