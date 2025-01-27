package wgetApp

import (
	"fmt"
	"net/http"
	"strings"

	wgetutils "wget/wgetUtils"

	"golang.org/x/net/html"
)

// mirror handles the mirroring of a webpage, downloading
// the page and recursively visits linked pages, downloading assets
// as needed.
func (app *WgetApp) mirror(url, rejectTypes, rejectPaths string, convertLink bool) error {
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
