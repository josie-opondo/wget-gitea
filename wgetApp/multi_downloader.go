package wgetApp

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	wgetutils "wget/wgetUtils"
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

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())

		if url == "" {
			continue // Skip empty lines
		}
		wg.Add(1)
		go func(url string) error {
			defer wg.Done()
			err := app.asyncDownload(outputFile, url, limit, directory)
			if err != nil {
				return err
			}
			return nil
		}(url)
	}
	wg.Wait()

	return nil
}

/*
asyncDownload
*parameters*
- outputFileName: The name of the file where the downloaded content will be saved. If empty, the filename is derived from the URL.
- url: The URL of the file to be downloaded.
- limit: The download speed limit (if applicable).
- directory: The directory where the file should be saved.

*functionality*
- Expands the directory path using `wgetutils.ExpandPath`.
- Sends an HTTP request to fetch the file.
- Validates the response status code to ensure a successful request.
- Determines the output file name, either from the provided `outputFileName` or derived from the URL.
- Creates the required directory structure if it does not exist.
- Opens a new file for writing the downloaded content.
- Reads the response body in chunks and writes to the file, applying a rate limit if specified.
- Displays download progress and prints a success message upon completion.
*/
func (app *WgetApp) asyncDownload(outputFileName, url, limit, directory string) error {
	path, err := wgetutils.ExpandPath(directory)
	if err != nil {
		return err
	}

	resp, err := wgetutils.HttpRequest(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: status %s url:\n[%s]", resp.Status, url)
	}

	if outputFileName == "" {
		urlParts := strings.Split(url, "/")
		fileName := urlParts[len(urlParts)-1]
		outputFileName = filepath.Join(path, fileName)
	} else {
		outputFileName = filepath.Join(path, outputFileName)
	}

	if path != "" {
		err = os.MkdirAll(path, 0o755)
		if err != nil {
			return fmt.Errorf("error creating directory:\n%v", err)
		}
	}

	var out *os.File
	out, err = os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating file:\n%v", err)
	}
	defer out.Close()

	var reader io.Reader = resp.Body
	if limit != "" {
		reader = wgetutils.NewRateLimitedReader(resp.Body, limit)
	}

	buffer := make([]byte, 32*1024)
	fmt.Printf("Downloading.... [%s]\n", url)
	var downloaded int64
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("oops! error reading response body")
		}

		if n > 0 {
			if _, err := out.Write(buffer[:n]); err != nil {
				return fmt.Errorf("error writing to file:\n%v", err)
			}
			downloaded += int64(n)
		}

		if err == io.EOF {
			break
		}
	}

	// endTime := time.Now()
	fmt.Printf("\033[32mDownloaded\033[0m [%s]\n", url)

	return nil
}
