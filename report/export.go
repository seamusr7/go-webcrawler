package report

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/seamusr7/go-webcrawler/crawler"
)

func ExportToCSV(w io.Writer, pages []crawler.PageInfo) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// CSV header
	headers := []string{
		"URL", "Referrer", "StatusCode", "Title", "Meta Description", "H1 Count", "Canonical",
		"Header Levels", "Internal Links", "External Links",
		"Has Main", "Has Nav", "Has Footer", "Has Header",
		"Inline Style Tags", "Inline Script Tags", "Structured Data Count",
		"Image Count", "Anchor Count",
		"SEO Fix Suggestions",
	}
	writer.Write(headers)

	// CSV content per page
	for _, p := range pages {
		headerLevels := ""
		for i, lvl := range p.HeaderLevels {
			headerLevels += strconv.Itoa(lvl)
			if i != len(p.HeaderLevels)-1 {
				headerLevels += ","
			}
		}

		suggestions := generateFixSuggestions(p)

		writer.Write([]string{
			p.URL,
			p.Referrer,
			strconv.Itoa(p.StatusCode),
			p.Title,
			p.Description,
			strconv.Itoa(p.H1Count),
			p.Canonical,
			headerLevels,
			strconv.Itoa(p.InternalLinks),
			strconv.Itoa(p.ExternalLinks),
			strconv.FormatBool(p.HasMain),
			strconv.FormatBool(p.HasNav),
			strconv.FormatBool(p.HasFooter),
			strconv.FormatBool(p.HasHeader),
			strconv.Itoa(p.InlineStyleTags),
			strconv.Itoa(p.InlineScriptTags),
			strconv.Itoa(p.StructuredDataCount),
			strconv.Itoa(len(p.Images)),
			strconv.Itoa(len(p.AnchorTexts)),
			suggestions,
		})
	}
}

// generateFixSuggestions returns human-friendly suggestions based on missing data
func generateFixSuggestions(p crawler.PageInfo) string {
	suggestions := []string{}

	if p.Title == "" {
		suggestions = append(suggestions, "Add a <title> tag.")
	}
	if p.Description == "" {
		suggestions = append(suggestions, "Add a meta description.")
	}
	if p.Canonical == "" {
		suggestions = append(suggestions, "Add a canonical tag.")
	}
	if p.H1Count == 0 {
		suggestions = append(suggestions, "Include at least one <h1> tag.")
	}
	if p.H1Count > 1 {
		suggestions = append(suggestions, "Reduce to one <h1> tag.")
	}
	if !p.HasMain {
		suggestions = append(suggestions, "Add a <main> tag for accessibility.")
	}
	if !p.HasNav {
		suggestions = append(suggestions, "Add a <nav> tag for navigation.")
	}
	if !p.HasFooter {
		suggestions = append(suggestions, "Add a <footer> tag.")
	}
	if !p.HasHeader {
		suggestions = append(suggestions, "Add a <header> tag.")
	}
	if p.StructuredDataCount == 0 {
		suggestions = append(suggestions, "Add structured data (ld+json).")
	}
	if len(p.Images) > 0 {
		missingAlt := 0
		for _, img := range p.Images {
			if img.Alt == "" {
				missingAlt++
			}
		}
		if missingAlt > 0 {
			suggestions = append(suggestions, strconv.Itoa(missingAlt)+" images are missing alt text.")
		}
	}

	return joinIfNotEmpty(suggestions, " | ")
}

// joinIfNotEmpty joins string slices if not empty
func joinIfNotEmpty(items []string, sep string) string {
	if len(items) == 0 {
		return "None"
	}
	return "\"" + join(items, sep) + "\""
}

// custom join for escaping commas inside values
func join(items []string, sep string) string {
	result := ""
	for i, item := range items {
		result += item
		if i < len(items)-1 {
			result += sep
		}
	}
	return result
}
