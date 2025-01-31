package wgetApp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
downloadMultipleFiles
*parameters*
- filePath: The path to the file containing URLs (one per line).
- outputFile: The output file where downloaded content is stored.
- limit: A concurrency limit for the number of simultaneous downloads.
- directory: The directory where files should be saved.

*functionality*
- Opens the file containing the URLs.
- Reads URLs line by line, skipping empty lines.
- Uses a WaitGroup to manage multiple concurrent downloads.
-  Calls AsyncDownload to handle each URL download asynchronously.
- Waits for all goroutines to complete before returning.
*/
func (app *WgetApp) downloadMultipleFiles(filePath, outputFile, limit, directory string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file:\n%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())

		if url == "" {
			continue // Skip empty lines
		}

		// Process each URL synchronously (one after the other)
		err := app.singleDownloader(outputFile, url, limit, directory)
		if err != nil {
			return err
		}
	}

	return nil
}
