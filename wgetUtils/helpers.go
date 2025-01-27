package wgetutils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ValidateURL checks if the given link is a valid URL.
func ValidateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	return nil
}

// ExtractDomain extracts the domain from a given URL string.
func ExtractDomain(urls string) (string, error) {
	u, err := url.Parse(urls)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// IsRejected checks if the URL ends with any of the rejected file types based on the provided rejectTypes string.
func IsRejected(url, rejectTypes string) bool {
	if rejectTypes == "" {
		return false
	}

	rejectedTypes := strings.Split(rejectTypes, ",")
	for _, ext := range rejectedTypes {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}
	return false
}

// HttpRequest sends an HTTP GET request to the provided URL with custom headers
// to simulate a browser request.
func HttpRequest(url string) (*http.Response, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request with a User-Agent header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers to mimic a Chrome browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, err
}

// isValidAttribute checks if an HTML tag attribute is valid for processing
func IsValidAttribute(tagName, attrKey string) bool {
	return (tagName == "link" && attrKey == "href") ||
		(tagName == "a" && attrKey == "href") ||
		(tagName == "script" && attrKey == "src") ||
		(tagName == "img" && attrKey == "src")
		
}

// ResolveURL resolves a relative URL to an absolute URL based on the given base URL.
// It handles fragment identifiers, protocols, and relative paths (e.g., './', '/', etc.).
func ResolveURL(base, rel string) string {
	// Remove fragment identifiers (anything starting with #)
	if fragmentIndex := strings.Index(rel, "#"); fragmentIndex != -1 {
		rel = rel[:fragmentIndex]
	}

	if strings.HasPrefix(rel, "http") {
		return rel
	}

	if strings.HasPrefix(rel, "//") {
		protocol := "http:"
		if strings.HasPrefix(base, "https") {
			protocol = "https:"
		}
		return protocol + rel
	}

	if strings.HasPrefix(rel, "/") {
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel
	}
	if strings.HasPrefix(rel, "./") {
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel[1:]
	}
	if strings.HasPrefix(rel, "//") && strings.Contains(rel[2:], "/") {
		baseParts := strings.Split(base, "/")
		return baseParts[0] + "//" + baseParts[2] + rel[1:]
	}

	baseParts := strings.Split(base, "/")
	return baseParts[0] + "//" + baseParts[2] + "/" + rel
}

// IsRejectedPath checks if the given URL contains any path specified in the pathRejects string.
func IsRejectedPath(url, pathRejects string) bool {
	if pathRejects == "" {
		return false
	}

	rejects := strings.Split(pathRejects, ",")
	for _, path := range rejects {
		if path[0] != '/' {
			continue
		}
		if contains(url, path) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a specified substring.
// It performs a simple substring search by comparing slices of the string.
func contains(str, substr string) bool {
	for i := 0; i < len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}