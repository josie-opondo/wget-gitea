package wgetutils

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

// ExpandPath expands shorthand notations to full paths
func ExpandPath(path string) (string, error) {
	// 1. Expand `~` to the home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error finding home directory:\n %v", err)
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// 2. Expand environment variables like $HOME, $USER, etc.
	path = os.ExpandEnv(path)

	// 3. Convert relative paths (./ or ../) to absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path:\n %v", err)
	}

	return absPath, nil
}

// FileExists checks if a file or directory exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// FormatSpeed formats a speed value into a human-readable format (GB, MB, or KiB).
func FormatSpeed(speed float64) string {
	if speed > 1000000000 {
		resGB := speed / 1000000000
		return fmt.Sprintf("~%.2fGB", roundToTwoDecimalPlaces(resGB))
	} else if speed > 1000000 {
		resMB := speed / 1000000
		return fmt.Sprintf("~%.2fMB", roundToTwoDecimalPlaces(resMB))
	}
	return fmt.Sprintf("%.0fKiB", speed)
}

// roundToTwoDecimalPlaces rounds a floating-point number to two decimal places.
func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

// LoadShowProgressState reads the showProgress state from a temporary file.
func LoadShowProgressState(tempConfigFile string) (bool, error) {
	if _, err := os.Stat(tempConfigFile); os.IsNotExist(err) {
		// File doesn't exist, return default true
		return true, nil
	}

	data, err := os.ReadFile(tempConfigFile)
	if err != nil {
		return false, fmt.Errorf("error reading showProgress state: %v", err)
	}

	// Parse the boolean value
	showProgress, err := strconv.ParseBool(string(data))
	if err != nil {
		return false, fmt.Errorf("error parsing showProgress state: %v", err)
	}

	// Delete the file after retrieving the state
	err = os.Remove(tempConfigFile)
	if err != nil {
		return false, fmt.Errorf("error deleting temp file: %v", err)
	}

	return showProgress, nil
}

// SaveProgressState saves the showProgress state to a temporary file.
func SaveProgressState(tempConfigFile string, showProgress bool) error {
	data := []byte(strconv.FormatBool(showProgress))
	err := os.WriteFile(tempConfigFile, data, 0o644)
	if err != nil {
		return fmt.Errorf("error saving showProgress state: %v", err)
	}
	return nil
}