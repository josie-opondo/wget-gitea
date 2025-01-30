package wgetApp

import "sync"

// UrlArgs struct with exported fields (Uppercase names)
type UrlArgs struct {
	url              string
	file             string
	rateLimit        string
	path             string
	sourceFile       string
	workInBackground bool
	mirroring        bool
	rejectFlag       string
	excludeFlag      string
	convertLinksFlag bool
}

// ProcessedURLs is a thread-safe structure that holds a collection of URLs
type ProcessedURLs struct {
	sync.Mutex
	urls map[string]bool
}

// WgetApp encapsulates global variables and synchronization primitives
type WgetApp struct {
	urlArgs        UrlArgs
	processedURLs  ProcessedURLs
	visitedPages   map[string]bool
	visitedAssets  map[string]bool
	muPages        sync.Mutex
	muAssets       sync.Mutex
	semaphore      chan struct{}
	count          int
	tempConfigFile string
}

// newWgetState initializes and returns a new instance of WgetApp.
func newWgetState() *WgetApp {
	return &WgetApp{
		visitedPages:  make(map[string]bool),
		visitedAssets: make(map[string]bool),
		processedURLs: ProcessedURLs{
			urls: make(map[string]bool),
		},
		semaphore: make(chan struct{}),
		count: 0,
		tempConfigFile: "progress_config.txt",
	}
}
