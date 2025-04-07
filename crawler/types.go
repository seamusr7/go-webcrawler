// types.go - Shared data types used across the SEO crawler

package crawler

import "time"

// ImageInfo holds details about <img> tags found on a page
type ImageInfo struct {
	Src        string // The image source URL
	Alt        string // The image alt text
	StatusCode int    // The HTTP status of the image (if checked)
}

// PageInfo holds all collected SEO-related data from a page
type PageInfo struct {
	URL                 string         // Page URL
	Referrer            string         // The URL that linked to this one
	StatusCode          int            // HTTP status code of the page
	Title               string         // <title> content
	Description         string         // <meta name="description">
	H1Count             int            // Count of <h1> tags
	Canonical           string         // <link rel="canonical">
	HeaderLevels        []int          // Order of heading levels (e.g. [1, 2, 3])
	Images              []ImageInfo    // List of image metadata
	AnchorTexts         []string       // Text content of <a> tags
	InternalLinks       int            // Count of internal links
	ExternalLinks       int            // Count of external links
	HasMain             bool           // Whether <main> tag exists
	HasNav              bool           // Whether <nav> tag exists
	HasFooter           bool           // Whether <footer> tag exists
	HasHeader           bool           // Whether <header> tag exists
	InlineStyleTags     int            // Number of inline <style> tags
	InlineScriptTags    int            // Number of inline <script> tags
	StructuredDataCount int            // Number of <script type="application/ld+json">
	Keywords            map[string]int // Word frequency map (for content/keyword analysis)
	MobileFriendly      bool           // True if viewport meta tag is detected
	CrawledAt           time.Time      // Timestamp of last crawl (for recrawling logic)
}
