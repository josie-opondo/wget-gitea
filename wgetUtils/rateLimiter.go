package wgetutils

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type RateLimitedReader struct {
	reader     io.Reader
	rateLimit  int64 // bytes per second
	bucket     int64
	lastFilled time.Time
}

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

func parseRateLimit(rateLimit string) (int64, error) {
	if len(rateLimit) < 2 {
		return 0, fmt.Errorf("invalid rate limit")
	}

	multiplier := 1
	switch rateLimit[len(rateLimit)-1] {
	case 'k', 'K':
		multiplier = 1024
		rateLimit = rateLimit[:len(rateLimit)-1]
	case 'M':
		multiplier = 1024 * 1024
		rateLimit = rateLimit[:len(rateLimit)-1]
	}

	rate, err := strconv.Atoi(rateLimit)
	if err != nil {
		return 0, err
	}
	return int64(rate * multiplier), nil
}

func NewRateLimitedReader(reader io.Reader, limit string) *RateLimitedReader {
	// Convert limit to bytes per second (rateLimit)
	rateLimit, _ := parseRateLimit(limit)
	return &RateLimitedReader{reader: reader, rateLimit: rateLimit, lastFilled: time.Now()}
}

func (r *RateLimitedReader) Read(p []byte) (n int, err error) {
	if r.bucket <= 0 {
		time.Sleep(time.Second)
		r.bucket = r.rateLimit
		r.lastFilled = time.Now()
	}

	toRead := int64(len(p))
	if toRead > r.bucket {
		toRead = r.bucket
	}

	n, err = r.reader.Read(p[:toRead])
	r.bucket -= int64(n)

	return n, err
}