package wgetApp

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	wgetutils "wget/wgetUtils"

	"golang.org/x/net/html"
)

// mirror fetches the content of a URL, processes it to extract links and assets,
// and downloads them if they belong to the same domain. It supports handling of
// inline styles and <style> tags, as well as recursive mirroring of linked pages.
// It also handles URL conversion for offline viewing if the convertLink flag is true.
func (app *WgetApp) mirror(url, rejectTypes, rejectPaths string, convertLink bool) error {
	app.muPages.Lock()
	if app.visitedPages[url] {
		app.muPages.Unlock()
		return nil
	}
	app.visitedPages[url] = true
	app.muPages.Unlock()
	
	domain, err := wgetutils.ExtractDomain(url)
	if err != nil {
		return fmt.Errorf("could not extract domain name for:\n%serror: %v", url, err)
	}

	// Check if we're at the root domain and force download of index.html
	if (strings.TrimRight(url, "/") == "http://"+domain || strings.TrimRight(url, "/") == "https://"+domain) && app.count == 0 {
		app.count++
		indexURL := strings.TrimRight(url, "/")
		app.downloadAsset(indexURL, domain, rejectTypes)
	}

	// Fetch and get the HTML of the page
	doc, err := fetchAndParsePage(url)
	if err != nil {
		return fmt.Errorf("error fetching or parsing page:\n%v", err)
	}

	// Function to handle links and assets found on the page
	handleLink := func(link, tagName string) {
		app.semaphore <- struct{}{}
		defer func() { <-app.semaphore }()

		baseURL := wgetutils.ResolveURL(url, link)
		if wgetutils.IsRejectedPath(baseURL, rejectPaths) {
			fmt.Printf("Skipping Rejected file path: %s\n", baseURL)
			return
		}
		baseURLDomain, err := wgetutils.ExtractDomain(baseURL)
		if err != nil {
			fmt.Println("Could not extract domain name for:", baseURLDomain, "\nError:", err)
			return
		}

		if baseURLDomain == domain {
			if tagName == "a" {
				if strings.HasSuffix(baseURL, "/") || strings.HasSuffix(baseURL, "/index.html") {
					// Ensure index.html is downloaded first
					indexURL := strings.TrimRight(baseURL, "/") + "/index.html"
					if !app.visitedPages[indexURL] {
						app.downloadAsset(indexURL, domain, rejectTypes)
						app.mirror(indexURL, rejectTypes, rejectPaths, convertLink)
					}
				} else {
					app.mirror(baseURL, rejectTypes, rejectPaths, convertLink)
				}
			}
			app.downloadAsset(baseURL, domain, rejectTypes)
		}
	}

	var wg sync.WaitGroup
	var processNode func(n *html.Node)

	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if wgetutils.IsValidAttribute(n.Data, attr.Key) {
					link := attr.Val
					if link != "" {
						wg.Add(1)
						go func(link, tagName string) {
							defer wg.Done()
							handleLink(link, tagName)
						}(link, n.Data)
					}
				}
				// Check for inline styles
				if attr.Key == "style" {
					app.extractAndHandleStyleURLs(attr.Val, url, domain, rejectTypes)
				}
			}
			// Check for <style> tags
			if n.Data == "style" && n.FirstChild != nil {
				app.extractAndHandleStyleURLs(n.FirstChild.Data, url, domain, rejectTypes)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	// Start processing the document
	processNode(doc)

	// Wait for all goroutines to complete
	wg.Wait()

	// Convert links if the flag is set
	if convertLink {
		wgetutils.ConvertLinks(url)
	}
	return nil
}

// downloadAsset checks if the asset URL has been visited, validates the URL, and initiates the download process.
func (app *WgetApp) downloadAsset(fileURL, domain, rejectTypes string) {
	app.muAssets.Lock()
	if app.visitedAssets[fileURL] {
		app.muAssets.Unlock()
		return
	}
	app.visitedAssets[fileURL] = true
	app.muAssets.Unlock()

	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("Invalid URL: %s\n", fileURL)
		return
	}

	if wgetutils.IsRejected(fileURL, rejectTypes) {
		fmt.Printf("Skipping rejected file: %s\n", fileURL)
		return
	}

	fmt.Printf("Downloading: %s\n", fileURL)
	app.asyncMirror("", fileURL, domain)
}

// fetchAndParsePage fetches the content of the URL and parses it as HTML
func fetchAndParsePage(url string) (*html.Node, error) {
	resp, err := wgetutils.HttpRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %s", resp.Status)
	}

	return html.Parse(resp.Body)
}

// extractAndHandleStyleURLs processes the URLs found in a CSS style block.
// It resolves relative URLs to absolute ones based on the base URL and downloads the assets,
// checking against domain restrictions and rejected types.
func (app *WgetApp) extractAndHandleStyleURLs(styleContent, baseURL, domain, rejectTypes string) {
	re := regexp.MustCompile(`url\(['"]?([^'"()]+)['"]?\)`)
	matches := re.FindAllStringSubmatch(styleContent, -1)
	for _, match := range matches {
		if len(match) > 1 {
			assetURL := wgetutils.ResolveURL(baseURL, match[1])
			app.downloadAsset(assetURL, domain, rejectTypes)
		}
	}
}
