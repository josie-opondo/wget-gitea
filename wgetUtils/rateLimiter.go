package wgetutils

import (
	"fmt"
	"strconv"
	"strings"
)

// RateLimitValidator validates the rate limit format for --rate-limit argument.
func RateLimitValidator(s string) error {
	ln := len(s) - 1
	idx := strings.Index(s, "=")
	if idx == -1 || !strings.ContainsAny(s[idx:], "kKmM") {
		return fmt.Errorf("invalid rate limit value.\nUsage: --rate-limit=400k || --rate-limit=2M")
	}

	// Extract the value and convert to integer
	val := s[idx+1 : ln]
	_, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("invalid rate limit value")
	}

	// Check that the rate limit ends with 'k' or 'm'
	if !strings.ContainsAny(s[ln:], "kKmM") {
		return fmt.Errorf("invalid rate limit value.\nUsage: --rate-limit=400k || --rate-limit=2M")
	}

	return nil
}
