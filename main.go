package main

import (
	"fmt"
	"os"

	"wget/wgetApp"
)

func main() {
	// Check if at least one argument (URL) is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <URL> [options]")
		return
	}

	// Initialize the WgetApp instance using the singleton pattern
	_, err := wgetApp.InitWget()
	if err != nil {
		// Print any initialization errors and exit
		fmt.Println(err)
		return
	}
}
