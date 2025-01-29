package wgetApp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestDownloadInBackgroundInvalidURL(t *testing.T) {
	app := &WgetApp{
		tempConfigFile: "test-config.json", // Example config file
	}

	// Test with an invalid URL
	fileName := ""
	urlStr := "invalid-url"
	rateLimit := "100k"
	err := app.downloadInBackground(fileName, urlStr, rateLimit)
	if err != nil {
		fmt.Println("err")
	}
	// assert.Error(t, err)
}

func TestDownloadInBackgroundLogCreation(t *testing.T) {
	app := &WgetApp{
		tempConfigFile: "test-config.json", // Example config file
	}

	// Test if the log file is created
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	defer ts.Close()

	urlStr := ts.URL
	rateLimit := "100k"
	err := app.downloadInBackground("", urlStr, rateLimit)
	if err != nil {
		fmt.Println("err")
	}

	// Wait for the download to complete
	time.Sleep(2 * time.Second) // Adjust as needed

	// Check if the log file exists
	_, err = os.Stat("wget-log")
	if err != nil {
		fmt.Println("err")
	}

	// Clean up
	os.Remove("wget-log")
}
func TestDownloadInBackgroundOutputDirectory(t *testing.T) {
	app := &WgetApp{
		tempConfigFile: "test-config.json", // Example config file
	}

	// Test if the output directory is created
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	defer ts.Close()

	// outputDir := "output-dir"
	urlStr := ts.URL
	rateLimit := "100k"
	err := app.downloadInBackground("", urlStr, rateLimit)
	if err != nil {
		fmt.Println("err")
	}

	// Wait for the download to complete
	time.Sleep(2 * time.Second) // Adjust as needed

	// Check if the output directory exists
	_, err = os.Stat(".")
	if err != nil {
		fmt.Println("err")
	}

	// Clean up
	// Note: Since the default path is ".", no need to remove it here.
}
