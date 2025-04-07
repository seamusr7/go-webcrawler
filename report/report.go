// report/report.go
package report

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/seamusr7/go-webcrawler/crawler"
)

// Generate generates a full SEO report for all crawled pages.
func Generate(pages []crawler.PageInfo) {
	fmt.Println("📋 SEO + Dead Link Report:")

	titleMap := make(map[string][]string)
	metaMap := make(map[string][]string)

	for _, page := range pages {
		fmt.Printf("- %s [%d] Title: %q | Meta: %q\n", page.URL, page.StatusCode, page.Title, page.Description)

		// STRUCTURAL SEMANTIC TAGS
		if !page.HasMain {
			fmt.Printf("  [SEO]   Missing <main> element on %s\n", page.URL)
		}
		if !page.HasNav {
			fmt.Printf("  [SEO]   Missing <nav> element on %s\n", page.URL)
		}
		if !page.HasFooter {
			fmt.Printf("  [SEO]   Missing <footer> element on %s\n", page.URL)
		}
		if !page.HasHeader {
			fmt.Printf("  [SEO]   Missing <header> element on %s\n", page.URL)
		}

		// ANCHOR TEXT ANALYSIS
		emptyAnchors := 0
		for _, text := range page.AnchorTexts {
			if strings.TrimSpace(text) == "" {
				emptyAnchors++
			}
		}
		if emptyAnchors > 0 {
			fmt.Printf("  [SEO]   %d anchor tags missing link text on %s\n", emptyAnchors, page.URL)
		}

		// LINK METRICS
		fmt.Printf("  [Links] Internal: %d | External: %d\n", page.InternalLinks, page.ExternalLinks)

		// IMAGE ANALYSIS
		if len(page.Images) > 0 {
			missingAlt := 0
			altMap := make(map[string]int)
			for _, img := range page.Images {
				if img.Alt == "" {
					missingAlt++
				} else {
					altMap[img.Alt]++
				}
			}
			if missingAlt > 0 {
				fmt.Printf("  [SEO]   %d image(s) missing alt text\n", missingAlt)
			}
			for alt, count := range altMap {
				if count > 1 {
					fmt.Printf("  [SEO]   Duplicate alt text %q used %d times\n", alt, count)
				}
			}
		}

		// BASIC TAG CHECKS
		if page.StatusCode >= 400 {
			fmt.Printf("  [ERROR] Broken link (%d): %s → %s\n", page.StatusCode, page.Referrer, page.URL)
		}
		if page.Title == "" {
			fmt.Printf("  [SEO]   Missing <title> on %s\n", page.URL)
		}
		if page.Description == "" {
			fmt.Printf("  [SEO]   Missing meta description on %s\n", page.URL)
		}
		if page.H1Count == 0 {
			fmt.Printf("  [SEO]   Missing <h1> on %s\n", page.URL)
		}
		if page.H1Count > 1 {
			fmt.Printf("  [SEO]   Multiple <h1> tags (%d) on %s\n", page.H1Count, page.URL)
		}
		if page.Canonical == "" {
			fmt.Printf("  [SEO]   Missing canonical tag on %s\n", page.URL)
		}

		// HEADING NESTING VALIDATION
		for i := 1; i < len(page.HeaderLevels); i++ {
			curr := page.HeaderLevels[i]
			prev := page.HeaderLevels[i-1]
			if curr > prev+1 {
				fmt.Printf("  [SEO]   Skipped heading level from <h%d> to <h%d> on %s\n", prev, curr, page.URL)
				break
			}
		}

		// INLINE STYLE/SCRIPT WARNINGS
		if page.InlineStyleTags > 5 {
			fmt.Printf("  [SEO]   High number of inline <style> tags: %d on %s\n", page.InlineStyleTags, page.URL)
		}
		if page.InlineScriptTags > 5 {
			fmt.Printf("  [SEO]   High number of inline <script> tags: %d on %s\n", page.InlineScriptTags, page.URL)
		}
		if page.StructuredDataCount == 0 {
			fmt.Printf("  [SEO]   No structured data (ld+json) found on %s\n", page.URL)
		}

		// KEYWORD DENSITY CHECK (on title and description)
		titleWords := extractWords(page.Title)
		descriptionWords := extractWords(page.Description)
		wordFreq := make(map[string]int)
		for _, word := range append(titleWords, descriptionWords...) {
			wordFreq[strings.ToLower(word)]++
		}
		for word, count := range wordFreq {
			if count > 2 && len(word) > 3 {
				fmt.Printf("  [SEO]   Overused keyword %q (%d times) in title/meta on %s\n", word, count, page.URL)
			}
		}

		// TRACK FOR DUPLICATES
		if page.Title != "" {
			titleMap[page.Title] = append(titleMap[page.Title], page.URL)
		}
		if page.Description != "" {
			metaMap[page.Description] = append(metaMap[page.Description], page.URL)
		}
	}

	// DUPLICATE CONTENT REPORT
	fmt.Println("\n🔁 Duplicate SEO Content Check:")
	for title, urls := range titleMap {
		if len(urls) > 1 {
			fmt.Printf("  [SEO]   Duplicate title: %q used on:\n", title)
			for _, u := range urls {
				fmt.Printf("           - %s\n", u)
			}
		}
	}
	for meta, urls := range metaMap {
		if len(urls) > 1 {
			fmt.Printf("  [SEO]   Duplicate meta description: %q used on:\n", meta)
			for _, u := range urls {
				fmt.Printf("           - %s\n", u)
			}
		}
	}
}

// extractWords splits a string into cleaned words for analysis
func extractWords(s string) []string {
	var words []string
	curr := ""
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			curr += string(r)
		} else if curr != "" {
			words = append(words, curr)
			curr = ""
		}
	}
	if curr != "" {
		words = append(words, curr)
	}
	return words
}
