package wgetApp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestWgetApp_asyncMirror(t *testing.T) {
	// Step 1: Create a temporary directory for testing
	tempDir := t.TempDir()
	testFileName := "testfile.txt"
	testURL := "/test/" + testFileName
	testContent := "This is a test file content."

	// Step 2: Start a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testContent))
	}))
	defer server.Close()

	// Step 3: Create a mock WgetApp that matches your model
	app := &WgetApp{
		urlArgs:       UrlArgs{},
		visitedPages:  make(map[string]bool),
		visitedAssets: make(map[string]bool),
		processedURLs: ProcessedURLs{urls: make(map[string]bool)}, // No pointer here!
	}

	// Step 4: Run asyncMirror
	err := app.asyncMirror("", server.URL+testURL, tempDir)
	if err != nil {
		t.Fatalf("asyncMirror failed: %v", err)
	}

	
}
