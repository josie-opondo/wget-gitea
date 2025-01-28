package wgetApp

import "sync"

type AppState struct {
	urlArgs UrlArgs
	// rateLimitedReader RateLimitedReader
	processedURLs  ProcessedURLs
	visitedPages   map[string]bool
	visitedAssets  map[string]bool
	muPages        sync.Mutex
	muAssets       sync.Mutex
	semaphore      chan struct{}
	count          int
	tempConfigFile string
}