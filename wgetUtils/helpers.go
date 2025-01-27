package wgetutils

import (
	"fmt"
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
