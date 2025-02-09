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

	// Step 5: Verify the file exists
	expectedFilePath := filepath.Join(tempDir, "test", testFileName)
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected file %s not created", expectedFilePath)
	}

	// Step 6: Read and verify file contents
	data, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(data) != testContent {
		t.Errorf("File content mismatch: expected %q, got %q", testContent, string(data))
	}

	// Step 7: Ensure duplicate URL is handled correctly
	err = app.asyncMirror("", server.URL+testURL, tempDir)
	if err == nil {
		t.Errorf("Expected error for duplicate URL, but got nil")
	}
}
