package wgetutils

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// convertCSSURLs replaces all URL references in a CSS file with local file system paths.
// This ensures that external assets referenced in stylesheets are properly mapped for offline use.
func convertCSSURLs(cssContent string) string {
	re := regexp.MustCompile(`url\(([^)]+)\)`)
	return re.ReplaceAllStringFunc(cssContent, func(match string) string {
		url := strings.Trim(match[4:len(match)-1], "'\"")
		localPath := getLocalPath(url)
		return fmt.Sprintf("url('%s')", localPath)
	})
}

// getLocalPath converts a given URL into a local file system path.
// It handles absolute HTTP(S) URLs, protocol-relative URLs, and root-relative paths.
func getLocalPath(originalURL string) string {
	if strings.HasPrefix(originalURL, "http") || strings.HasPrefix(originalURL, "//") {
		parsedURL, err := url.Parse(originalURL)
		if err != nil {
			return originalURL
		}
		return path.Join(parsedURL.Host, parsedURL.Path)
	} else if strings.HasPrefix(originalURL, "/") {
		return path.Join(".", originalURL)
	}
	return originalURL
}

// removeHTTP removes the http:// or https:// prefix from the URL.
func removeHTTP(url string) string {
	re := regexp.MustCompile(`^https?://`)
	modifiedURL := re.ReplaceAllString(url, "")
	isBaseURL := regexp.MustCompile(`^[^/]+/?$`).MatchString(modifiedURL)

	// If the URL is a base URL, append "index.html" if it's not already present
	if isBaseURL {
		if strings.HasSuffix(modifiedURL, "/") {
			modifiedURL += "index.html"
		} else {
			modifiedURL += "/index.html"
		}
	}

	return modifiedURL
}
