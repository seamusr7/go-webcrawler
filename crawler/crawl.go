package crawler

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// Crawl fetches a webpage, parses its DOM, and extracts SEO-relevant data.
func Crawl(currentURL, referrer string) (PageInfo, []string, error) {
	fmt.Println("✨ Crawling:", currentURL)

	resp, err := http.Get(currentURL)
	if err != nil {
		return PageInfo{URL: currentURL, Referrer: referrer, StatusCode: 0}, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return PageInfo{URL: currentURL, Referrer: referrer, StatusCode: statusCode}, nil, err
	}

	base, _ := url.Parse(currentURL)
	var links []string
	var title, description, canonical string
	var h1Count int
	var headerLevels []int
	var images []ImageInfo
	var anchorTexts []string
	var internalLinks, externalLinks int
	var hasMain, hasNav, hasFooter, hasHeader bool
	var inlineStyleTags, inlineScriptTags, structuredDataCount int

	// Traverse DOM
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href, err := url.Parse(attr.Val)
						if err == nil {
							absolute := base.ResolveReference(href)
							links = append(links, absolute.String())
							if absolute.Host == base.Host {
								internalLinks++
							} else {
								externalLinks++
							}
						}
					}
				}
				if n.FirstChild != nil {
					anchorTexts = append(anchorTexts, n.FirstChild.Data)
				}

			case "title":
				if n.FirstChild != nil {
					title = n.FirstChild.Data
				}

			case "meta":
				var nameAttr, contentAttr string
				for _, attr := range n.Attr {
					if attr.Key == "name" && attr.Val == "description" {
						nameAttr = attr.Val
					}
					if attr.Key == "content" {
						contentAttr = attr.Val
					}
				}
				if nameAttr == "description" {
					description = contentAttr
				}

			case "link":
				for _, attr := range n.Attr {
					if attr.Key == "rel" && attr.Val == "canonical" {
						for _, a := range n.Attr {
							if a.Key == "href" {
								canonical = a.Val
							}
						}
					}
				}

			case "h1":
				h1Count++
				headerLevels = append(headerLevels, 1)
			case "h2":
				headerLevels = append(headerLevels, 2)
			case "h3":
				headerLevels = append(headerLevels, 3)
			case "h4":
				headerLevels = append(headerLevels, 4)
			case "h5":
				headerLevels = append(headerLevels, 5)
			case "h6":
				headerLevels = append(headerLevels, 6)

			case "img":
				var src, alt string
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						src = attr.Val
					}
					if attr.Key == "alt" {
						alt = attr.Val
					}
				}
				if src != "" {
					images = append(images, ImageInfo{Src: src, Alt: alt})
				}

			case "main":
				hasMain = true
			case "nav":
				hasNav = true
			case "footer":
				hasFooter = true
			case "header":
				hasHeader = true

			case "style":
				inlineStyleTags++

			case "script":
				isLDJSON := false
				for _, attr := range n.Attr {
					if attr.Key == "type" && attr.Val == "application/ld+json" {
						isLDJSON = true
					}
				}
				if isLDJSON {
					structuredDataCount++
				} else if len(n.Attr) == 0 && n.FirstChild != nil {
					inlineScriptTags++
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)

	// Construct result
	return PageInfo{
		URL:                 currentURL,
		Referrer:            referrer,
		StatusCode:          statusCode,
		Title:               title,
		Description:         description,
		H1Count:             h1Count,
		Canonical:           canonical,
		HeaderLevels:        headerLevels,
		Images:              images,
		AnchorTexts:         anchorTexts,
		InternalLinks:       internalLinks,
		ExternalLinks:       externalLinks,
		HasMain:             hasMain,
		HasNav:              hasNav,
		HasFooter:           hasFooter,
		HasHeader:           hasHeader,
		InlineStyleTags:     inlineStyleTags,
		InlineScriptTags:    inlineScriptTags,
		StructuredDataCount: structuredDataCount,
	}, links, nil
}
