// worker.go - Worker logic for handling crawl jobs
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

func StartCrawling(startURL string, maxPages int) []PageInfo {
	var (
		visited   = make(map[string]bool)
		visitedMu sync.Mutex
	)

	jobChan := make(chan [2]string, 100)
	pageResults := make(chan PageInfo, 100)
	linkResults := make(chan []string, 100)
	done := make(chan struct{})

	var wg sync.WaitGroup
	var managerWg sync.WaitGroup
	var allPages []PageInfo
	var resultsMu sync.Mutex

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
						jobChan <- [2]string{link, ""}
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

	visitedMu.Lock()
	visited[startURL] = true
	visitedMu.Unlock()
	wg.Add(1)
	jobChan <- [2]string{startURL, ""}

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go Worker(i, jobChan, pageResults, linkResults, &wg)
	}

	wg.Wait()

	close(done)
	managerWg.Wait()
	close(jobChan)
	close(linkResults)
	close(pageResults)

	fmt.Println("\n✅ Done crawling.")
	return allPages
}
