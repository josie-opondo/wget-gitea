package wgetApp

import "sync"

// UrlArgs struct with exported fields (Uppercase names)

// AppState encapsulates global variables and synchronization primitives
type AppState struct {
	urlArgs           UrlArgs
	// rateLimitedReader RateLimitedReader
	processedURLs     ProcessedURLs
	visitedPages      map[string]bool
	visitedAssets     map[string]bool
	muPages           sync.Mutex
	muAssets          sync.Mutex
	semaphore         chan struct{}
	count             int
	tempConfigFile    string
}