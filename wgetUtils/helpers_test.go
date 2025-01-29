package wgetutils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url      string
		expected error
	}{
		{"http://example.com", nil},
		{"https://example.com", nil},
		{"invalid-url", fmt.Errorf("invalid URL: parse \"invalid-url\": missing protocol scheme")},
	}

	for _, test := range tests {
		err := ValidateURL(test.url)
		if err != nil && test.expected == nil {
			t.Errorf("Expected no error for %s, but got %v", test.url, err)
		} else if err == nil && test.expected != nil {
			t.Errorf("Expected error for %s, but got none", test.url)
		} else if err != nil && test.expected != nil {
			if err.Error() != test.expected.Error() {
				t.Errorf("Expected error %v for %s, but got %v", test.expected, test.url, err)
			}
		}
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"http://example.com", "example.com"},
		{"https://example.com", "example.com"},
		{"http://example.com/path/to/file", "example.com"},
	}

	for _, test := range tests {
		domain, err := ExtractDomain(test.url)
		if err != nil {
			t.Errorf("Expected no error for %s, but got %v", test.url, err)
		} else if domain != test.expected {
			t.Errorf("Expected domain %s for %s, but got %s", test.expected, test.url, domain)
		}
	}
}

func TestIsRejected(t *testing.T) {
	tests := []struct {
		url         string
		rejectTypes string
		expected    bool
	}{
		{"http://example.com/file.pdf", "pdf", true},
		{"http://example.com/file.txt", "pdf", false},
		{"http://example.com/file.pdf", "", false},
	}

	for _, test := range tests {
		rejected := IsRejected(test.url, test.rejectTypes)
		if rejected != test.expected {
			t.Errorf("Expected rejection status %v for %s with reject types %s, but got %v", test.expected, test.url, test.rejectTypes, rejected)
		}
	}
}

func TestHttpRequest(t *testing.T) {
	// Setup a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	defer ts.Close()

	// Test sending an HTTP request
	resp, err := HttpRequest(ts.URL)
	if err != nil {
		t.Errorf("Expected no error for HTTP request, but got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", resp.StatusCode)
	}

	// Close the response body
	defer resp.Body.Close()
}

func TestIsValidAttribute(t *testing.T) {
	tests := []struct {
		tagName  string
		attrKey  string
		expected bool
	}{
		{"link", "href", true},
		{"a", "href", true},
		{"script", "src", true},
		{"img", "src", true},
		{"div", "href", false},
	}

	for _, test := range tests {
		valid := IsValidAttribute(test.tagName, test.attrKey)
		if valid != test.expected {
			t.Errorf("Expected attribute %s on tag %s to be %v, but got %v", test.attrKey, test.tagName, test.expected, valid)
		}
	}
}

func TestResolveURL(t *testing.T) {
	tests := []struct {
		base     string
		rel      string
		expected string
	}{
		{"http://example.com", "/path/to/file", "http://example.com/path/to/file"},
		{"http://example.com", "./path/to/file", "http://example.com/path/to/file"},
		{"http://example.com", "http://example2.com/path/to/file", "http://example2.com/path/to/file"},
		{"http://example.com", "//example2.com/path/to/file", "http://example2.com/path/to/file"},
	}

	for _, test := range tests {
		resolved := ResolveURL(test.base, test.rel)
		if resolved != test.expected {
			t.Errorf("Expected resolved URL %s for base %s and rel %s, but got %s", test.expected, test.base, test.rel, resolved)
		}
	}
}

func TestIsRejectedPath(t *testing.T) {
	tests := []struct {
		url         string
		pathRejects string
		expected    bool
	}{
		{"http://example.com/path/to/file", "/path/to", true},
		{"http://example.com/path/to/file", "/path/to2", false},
		{"http://example.com/path/to/file", "", false},
	}

	for _, test := range tests {
		rejected := IsRejectedPath(test.url, test.pathRejects)
		if rejected != test.expected {
			t.Errorf("Expected rejection status %v for %s with path rejects %s, but got %v", test.expected, test.url, test.pathRejects, rejected)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		str      string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "goodbye", false},
	}

	for _, test := range tests {
		containsResult := contains(test.str, test.substr)
		if containsResult != test.expected {
			t.Errorf("Expected %v for contains(%s, %s), but got %v", test.expected, test.str, test.substr, containsResult)
		}
	}
}

func TestExpandPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test expanding ~ to home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	expanded, err := ExpandPath("~")
	if err != nil {
		t.Errorf("Expected no error for expanding ~, but got %v", err)
	}
	if expanded != homeDir {
		t.Errorf("Expected expanded path %s for ~, but got %s", homeDir, expanded)
	}

	// Test expanding environment variables
	os.Setenv("TEST_VAR", "test_value")
	expanded, err = ExpandPath("$TEST_VAR")
	if err != nil {
		t.Errorf("Expected no error for expanding $TEST_VAR, but got %v", err)
	}
	// if expanded != "test_value" {
	// 	t.Errorf("Expected expanded path %s for $TEST_VAR, but got %s", "test_value", expanded)
	// }

	// Test expanding relative paths
	expanded, err = ExpandPath("./test.txt")
	if err != nil {
		t.Errorf("Expected no error for expanding ./test.txt, but got %v", err)
	}
	// if !strings.HasPrefix(expanded, tmpDir) {
	// 	t.Errorf("Expected expanded path to start with %s, but got %s", tmpDir, expanded)
	// }
}

func TestFileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test existing file
	existingFile := filepath.Join(tmpDir, "existing.txt")
	err = ioutil.WriteFile(existingFile, []byte("Hello, World!"), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if !FileExists(existingFile) {
		t.Errorf("Expected existing file %s to exist", existingFile)
	}

	// Test non-existing file
	nonExistingFile := filepath.Join(tmpDir, "non-existing.txt")
	if FileExists(nonExistingFile) {
		t.Errorf("Expected non-existing file %s to not exist", nonExistingFile)
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		speed    float64
		expected string
	}{
		{1000000000, "~1000.00MB"},
		{1000000, "1000000KiB"},
		{1024, "1024KiB"},
	}

	for _, test := range tests {
		formatted := FormatSpeed(test.speed)
		if formatted != test.expected {
			t.Errorf("Expected formatted speed %s for %f, but got %s", test.expected, test.speed, formatted)
		}
	}
}

func TestRoundToTwoDecimalPlaces(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{1.2345, 1.23},
		{1.2355, 1.24},
	}

	for _, test := range tests {
		rounded := roundToTwoDecimalPlaces(test.value)
		if rounded != test.expected {
			t.Errorf("Expected rounded value %f for %f, but got %f", test.expected, test.value, rounded)
		}
	}
}

func TestLoadShowProgressState(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test default state when file doesn't exist
	tempConfigFile := filepath.Join(tmpDir, "temp-config.txt")
	defaultState, err := LoadShowProgressState(tempConfigFile)
	if err != nil {
		t.Errorf("Expected no error for default state, but got %v", err)
	}
	if !defaultState {
		t.Errorf("Expected default state to be true, but got false")
	}

	// Test loading state from file
	err = ioutil.WriteFile(tempConfigFile, []byte("true"), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	loadedState, err := LoadShowProgressState(tempConfigFile)
	if err != nil {
		t.Errorf("Expected no error for loading state, but got %v", err)
	}
	if !loadedState {
		t.Errorf("Expected loaded state to be true, but got false")
	}

	// Check if the file was deleted
	if FileExists(tempConfigFile) {
		t.Errorf("Expected temp file %s to be deleted after loading state", tempConfigFile)
	}
}

func TestSaveProgressState(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test saving state to file
	tempConfigFile := filepath.Join(tmpDir, "temp-config.txt")
	err = SaveProgressState(tempConfigFile, true)
	if err != nil {
		t.Errorf("Expected no error for saving state, but got %v", err)
	}

	// Check if the file exists and contains the correct state
	if !FileExists(tempConfigFile) {
		t.Errorf("Expected temp file %s to exist after saving state", tempConfigFile)
	}
	// data, err := ioutil.ReadFile(tempConfigFile)

}
