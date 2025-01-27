package wgetutils

import (
	"fmt"
	"net/url"
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