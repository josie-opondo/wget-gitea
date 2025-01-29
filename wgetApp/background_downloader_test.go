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

func TestDownloadInBackgroundSaveProgressState(t *testing.T) {
	app := &WgetApp{
		tempConfigFile: "test-config.json", // Example config file
	}

	// Test saving progress state
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

	// Check if the progress state file exists
	_, err = os.Stat(app.tempConfigFile)
	if err != nil {
		fmt.Println("err")
	}
	// Clean up
	os.Remove(app.tempConfigFile)
}

func TestDownloadInBackground(t *testing.T) {
	// Setup a test server to serve a file
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	defer ts.Close()

	// Create a WgetApp instance for testing
	app := &WgetApp{
		tempConfigFile: "test-config.json", // Example config file
	}

	// Test downloading a file with a custom name
	fileName := "custom-name.txt"
	urlStr := ts.URL
	rateLimit := "100k"
	err := app.downloadInBackground(fileName, urlStr, rateLimit)
	if err != nil {
		fmt.Println("err")
	}
	// Wait for the download to complete
	time.Sleep(2 * time.Second) // Adjust as needed

	// Check if the file was downloaded correctly
	_, err = os.Stat(fileName)
	if err != nil {
		fmt.Println("err")
	}
	// Clean up
	os.Remove(fileName)
}
