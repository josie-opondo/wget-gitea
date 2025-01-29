package wgetApp

import (
	"fmt"
	"testing"
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
