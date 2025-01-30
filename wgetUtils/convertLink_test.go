package wgetutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestConvertLinks(t *testing.T) {
	// Create a test HTML file
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	htmlFilePath := filepath.Join(tmpDir, "test.html")
	err = ioutil.WriteFile(htmlFilePath, []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<style>
				body {
					background-image: url('http://example.com/background.jpg');
				}
			</style>
		</head>
		<body>
			<a href="http://example.com/page1.html">Page 1</a>
			<img src="http://example.com/image.jpg" />
		</body>
		</html>
	`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Run ConvertLinks
	ConvertLinks(htmlFilePath)

	// Read the modified HTML file
	modifiedHTML, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Check if links were converted correctly
	if !strings.Contains(string(modifiedHTML), "example.com/background.jpg") {
		t.Errorf("Expected background.jpg link to be converted")
	}
	if !strings.Contains(string(modifiedHTML), "example.com/page1.html") {
		t.Errorf("Expected page1.html link to be converted")
	}
	if !strings.Contains(string(modifiedHTML), "example.com/image.jpg") {
		t.Errorf("Expected image.jpg link to be converted")
	}
}

func TestModifyLinks(t *testing.T) {
	// Create a test HTML document
	doc, err := html.Parse(strings.NewReader(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<style>
				body {
					background-image: url('http://example.com/background.jpg');
				}
			</style>
		</head>
		<body>
			<a href="http://example.com/page1.html">Page 1</a>
			<img src="http://example.com/image.jpg" />
		</body>
		</html>
	`))
	if err != nil {
		t.Fatal(err)
	}

	// Modify the links
	modifyLinks(doc, ".")

	// Check if links were modified correctly
	var foundLinks []string
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					foundLinks = append(foundLinks, attr.Val)
				} else if attr.Key == "style" {
					re := regexp.MustCompile(`url\(([^)]+)\)`)
					matches := re.FindAllStringSubmatch(attr.Val, -1)
					for _, match := range matches {
						foundLinks = append(foundLinks, strings.Trim(match[1], "'\""))
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	// Check if all links were converted
	for _, link := range foundLinks {
		if strings.HasPrefix(link, "http") {
			t.Errorf("Expected link %s to be converted", link)
		}
	}
}

func TestConvertCSSURLs(t *testing.T) {
	// Test converting CSS URLs
	cssContent := "body { background-image: url('http://example.com/background.jpg'); }"
	modifiedCSS := convertCSSURLs(cssContent)

	// Check if the URL was converted correctly
	if !strings.Contains(modifiedCSS, "example.com/background.jpg") {
		t.Errorf("Expected background.jpg URL to be converted")
	}
}

func TestGetLocalPath(t *testing.T) {
	// Test converting URLs to local paths
	tests := []struct {
		url      string
		expected string
	}{
		{"http://example.com/path/to/file", "example.com/path/to/file"},
		{"//example.com/path/to/file", "example.com/path/to/file"},
		{"/path/to/file", "path/to/file"},
		{"relative/path/to/file", "relative/path/to/file"},
	}

	for _, test := range tests {
		localPath := getLocalPath(test.url)
		if localPath != test.expected {
			t.Errorf("Expected %s to be converted to %s, but got %s", test.url, test.expected, localPath)
		}
	}
}

func TestRemoveHTTP(t *testing.T) {
	// Test removing HTTP(S) prefixes from URLs
	tests := []struct {
		url      string
		expected string
	}{
		{"http://example.com", "example.com/index.html"},
		{"https://example.com", "example.com/index.html"},
		{"http://example.com/path/to/file", "example.com/path/to/file"},
		{"https://example.com/path/to/file", "example.com/path/to/file"},
		{"//example.com/path/to/file", "//example.com/path/to/file"},
		{"/path/to/file", "/path/to/file"},
		{"relative/path/to/file", "relative/path/to/file"},
	}

	for _, test := range tests {
		modifiedURL := removeHTTP(test.url)
		if modifiedURL != test.expected {
			t.Errorf("Expected %s to be converted to %s, but got %s", test.url, test.expected, modifiedURL)
		}
	}

	// Test appending "index.html" to base URLs
	tests = []struct {
		url      string
		expected string
	}{
		{"http://example.com", "example.com/index.html"},
		{"https://example.com", "example.com/index.html"},
		{"http://example.com/", "example.com/index.html"},
		{"https://example.com/", "example.com/index.html"},
	}

	for _, test := range tests {
		modifiedURL := removeHTTP(test.url)
		if modifiedURL != test.expected {
			t.Errorf("Expected %s to be converted to %s, but got %s", test.url, test.expected, modifiedURL)
		}
	}
}
