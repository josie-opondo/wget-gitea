package wgetApp

import "sync"

// Global variables for singleton pattern
var (
	state *WgetApp   // Holds the single instance of WgetApp
	once  sync.Once  // Ensures one-time initialization
)

// InitWget initializes and returns the singleton instance of WgetApp.
// It ensures that only one instance of WgetApp is created using sync.Once.
func InitWget() (*WgetApp, error) {
	var err error

	once.Do(func() {
		// Initialize the singleton instance
		state = newWgetState()
	})

	// Return any initialization errors (though none are expected here)
	if err != nil {
		return nil, err
	}

	return state, nil
}
