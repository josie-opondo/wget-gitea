package wgetutils

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestRateLimitValidator(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr bool
	}{
		{"--rate-limit=400k", false},
		{"--rate-limit=2M", false},
		{"--rate-limit=400", true},
		{"--rate-limit=2", true},
		{"--rate-limit=abc", true},
		{"--rate-limit=400K", false},
		{"--rate-limit=2M", false},
	}

	for _, test := range tests {
		err := RateLimitValidator(test.input)
		if err != nil && !test.expectedErr {
			t.Errorf("Expected no error for %s, but got %v", test.input, err)
		} else if err == nil && test.expectedErr {
			t.Errorf("Expected error for %s, but got none", test.input)
		}
	}
}

func TestParseRateLimit(t *testing.T) {
	tests := []struct {
		input       string
		expected    int64
		expectedErr bool
	}{
		{"400k", 400 * 1024, false},
		{"2M", 2 * 1024 * 1024, false},
		{"400", 400, false},
		{"abc", 0, true},
		{"400K", 400 * 1024, false},
		{"2M", 2 * 1024 * 1024, false},
	}

	for _, test := range tests {
		rateLimit, err := parseRateLimit(test.input)
		if err != nil && !test.expectedErr {
			t.Errorf("Expected no error for %s, but got %v", test.input, err)
		} else if err == nil && test.expectedErr {
			t.Errorf("Expected error for %s, but got none", test.input)
		} else if rateLimit != test.expected {
			t.Errorf("Expected rate limit %d for %s, but got %d", test.expected, test.input, rateLimit)
		}
	}
}

func TestNewRateLimitedReader(t *testing.T) {
	// Create a test reader
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	// defer os.Remove(tmpFile)

	// Write some data to the file
	_, err = tmpFile.WriteString("Hello, World!")
	if err != nil {
		t.Fatal(err)
	}
	// tmpFile.Close()

	// Open the file for reading
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	// defer file.Close()

	// Create a rate-limited reader
	rateLimitedReader := NewRateLimitedReader(file, "1k")

	// Check if the reader is correctly initialized
	if rateLimitedReader.rateLimit != 1024 {
		t.Errorf("Expected rate limit 1024, but got %d", rateLimitedReader.rateLimit)
	}
}

func TestRateLimitedReaderRead(t *testing.T) {
	// Create a test reader
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	// defer tmpFile.Close()

	// Write some data to the file
	_, err = tmpFile.WriteString(strings.Repeat("a", 1024*1024)) // 1MB
	if err != nil {
		t.Fatal(err)
	}
	// tmpFile.Close()

	// Open the file for reading
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	// defer file.Close()

	// Create a rate-limited reader with a rate limit of 1KB/s
	rateLimitedReader := NewRateLimitedReader(file, "1k")

	// Read data from the rate-limited reader
	buf := make([]byte, 1024)
	startTime := time.Now()
	n, err := rateLimitedReader.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the read operation took approximately 1 second
	elapsedTime := time.Since(startTime)
	if elapsedTime < 900*time.Millisecond || elapsedTime > 1100*time.Millisecond {
		t.Errorf("Expected read operation to take approximately 1 second, but took %v", elapsedTime)
	}

	// Check if the correct amount of data was read
	if n != 1024 {
		t.Errorf("Expected to read 1024 bytes, but read %d", n)
	}
}
