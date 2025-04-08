package crawler

import (
	"fmt"
	"sync"
)

// Worker pulls jobs from the job channel, crawls the page, and sends results back
func Worker(id int, jobs <-chan [2]string, pageResults chan<- PageInfo, linkResults chan<- []string, wg *sync.WaitGroup) {
	for job := range jobs {
		url := job[0]
		referrer := job[1]

		pageInfo, links, err := Crawl(url, referrer)
		if err != nil {
			fmt.Printf("❌ Worker %d error: %s\n", id, err)
			wg.Done()
			continue
		}

		pageResults <- pageInfo
		linkResults <- links
		wg.Done()
	}
}

// StartCrawling initializes the crawl, manages workers, and collects results
func StartCrawling(startURL string, maxPages int) []PageInfo {
	var (
		visited     = make(map[string]bool)
		visitedMu   sync.Mutex
		resultsMu   sync.Mutex
		allPages    []PageInfo
		jobChan     = make(chan [2]string, 100)
		pageResults = make(chan PageInfo, 100)
		linkResults = make(chan []string, 100)
		wg          sync.WaitGroup
	)

	// Seed first job
	visited[startURL] = true
	wg.Add(1)
	jobChan <- [2]string{startURL, ""}

	// Start 5 workers
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go Worker(i, jobChan, pageResults, linkResults, &wg)
	}

	// Collector goroutine
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Process results as they come in
LOOP:
	for {
		select {
		case page := <-pageResults:
			resultsMu.Lock()
			if len(allPages) < maxPages {
				allPages = append(allPages, page)
			}
			resultsMu.Unlock()

		case links := <-linkResults:
			for _, link := range links {
				visitedMu.Lock()
				if !visited[link] && len(allPages) < maxPages {
					visited[link] = true
					wg.Add(1)
					jobChan <- [2]string{link, ""}
				}
				visitedMu.Unlock()
			}
		case <-done:
			break LOOP
		}
	}

	close(jobChan)
	// ⚠️ Do NOT close pageResults or linkResults — workers may still be writing!

	fmt.Println("\n✅ Done crawling.")
	return allPages
}
