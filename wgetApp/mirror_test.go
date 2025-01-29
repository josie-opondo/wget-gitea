package wgetApp

import (
	"sync"
	"testing"
)

func TestMirror(t *testing.T) {
	app := newWgetState()
	app.semaphore = make(chan struct{}, 1)

	t.Run("Valid URL", func(t *testing.T) {
		err := app.mirror("http://example.com", "", "", false)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		err := app.mirror("invalid_url", "", "", true)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDownloadAsset(t *testing.T) {
	app := newWgetState()
	app.muAssets = sync.Mutex{}

	t.Run("Valid Asset", func(t *testing.T) {
		app.downloadAsset("http://example.com/image.jpg", "example.com", "")
		if !app.visitedAssets["http://example.com/image.jpg"] {
			t.Errorf("Expected asset to be visited")
		}
	})

// 	t.Run("Invalid Asset", func(t *testing.T) {
// 		app.downloadAsset("", "example.com", "")
// 		if _, exists := app.visitedAssets[" "]; exists {
// 			t.Errorf("Expected asset to be ignored")
// 		}
// 	})
// }
