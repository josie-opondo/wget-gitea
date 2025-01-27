package wgetApp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	wgetutils "wget/wgetUtils"
)

// asyncMirror handles the asynchronous mirroring process. Currently, it returns nil,
// but the function is set up to be expanded for asynchronous mirroring operations in the future.
func (app *WgetApp) asyncMirror(outputFile, urls, direc string) error {
	app.processedURLs.Lock()
	if processed, exists := app.processedURLs.urls[urls]; exists && processed {
		app.processedURLs.Unlock()
		return fmt.Errorf("URL already processed:\n%s", urls)
	}
	app.processedURLs.Unlock()

	// Parse the URL to get the path components
	u, err := url.Parse(urls)
	if err != nil {
		return fmt.Errorf("error parsing URL:\n%v", err)
	}

	// Create the necessary directories based on the URL path
	rootPath, err := wgetutils.ExpandPath(direc)
	if err != nil {
		return err
	}

	pathComponents := strings.Split(strings.Trim(u.Path, "/"), "/")
	relativeDirPath := filepath.Join(pathComponents[:len(pathComponents)-1]...)
	fullDirPath := filepath.Join(rootPath, relativeDirPath)
	fileName := pathComponents[len(pathComponents)-1]

	resp, err := wgetutils.HttpRequest(urls)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: status %s\nurl: %s", resp.Status, urls)
	}

	contentType := resp.Header.Get("Content-Type")

	if outputFile == "" {
		if fileName == "" || strings.HasSuffix(urls, "/") {
			fileName = "index.html"
		} else if contentType == "text/html" && !strings.HasSuffix(fileName, ".html") {
			fileName += ".html"
		}
		outputFile = filepath.Join(fullDirPath, fileName)
	} else {
		if contentType == "text/html" && !strings.HasSuffix(outputFile, ".html") {
			outputFile += ".html"
		}
		outputFile = filepath.Join(fullDirPath, outputFile)
	}

	if fullDirPath != "" {
		if _, err := os.Stat(fullDirPath); os.IsNotExist(err) {
			err = os.MkdirAll(fullDirPath, 0o755)
			if err != nil {
				return fmt.Errorf("error creating path:\n%v", err)
			}
		}
	}

	var out *os.File
	out, err = os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file:\n%v", err)
	}
	defer out.Close()

	var reader io.Reader = resp.Body

	// Get the content length for the progress (if available)
	if length := resp.Header.Get("Content-Length"); length != "" {
		_, err = strconv.ParseInt(length, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing Content-Length:\n%v", err)
		}
	}

	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64

	// Download the file while showing progress
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading response body")
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

	fmt.Printf("\n\033[32mDownloaded [%s]\033[0m\n", urls)

	// Mark the URL as processed
	app.processedURLs.Lock()
	app.processedURLs.urls[urls] = true
	app.processedURLs.Unlock()
	return nil
}
