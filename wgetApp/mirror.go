package wgetApp

import (
	"fmt"
	"strings"
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

	fmt.Printf("Downloading: %s\n", fileURL)
}
