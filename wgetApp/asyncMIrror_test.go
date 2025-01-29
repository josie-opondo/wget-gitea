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

	
}
