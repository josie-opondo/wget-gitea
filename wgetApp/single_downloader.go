package wgetApp

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	wgetutils "wget/wgetUtils"
)

func (app *WgetApp) singleDownloader(file, url, limit, directory string) error {
	path, err := wgetutils.ExpandPath(directory)
	if err != nil {
		return err
	}
	fileURL := url
	startTime := time.Now()
	toDisplay, err := wgetutils.LoadShowProgressState(app.tempConfigFile)
	if err != nil {
		return err
	}
	fmt.Printf("started at %s\n", startTime.Format("2006-01-02 15:04:05"))

	resp, err := wgetutils.HttpRequest(fileURL)
	if err != nil {
		return fmt.Errorf("error downloading file:\nserver misbehaving")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: status %s\nurl: [%s]", resp.Status, url)
	}
	fmt.Printf("sending request, awaiting response... status %s\n", resp.Status)

	contentLength := resp.ContentLength
	fmt.Printf("content size: %d bytes [~%.2fMB]\n", contentLength, float64(contentLength)/1000000)

	// Set the output file name
	var outputFile string
	if file == "" {
		urlParts := strings.Split(fileURL, "/")
		file = urlParts[len(urlParts)-1]
		outputFile = filepath.Join(path, file)
	} else {
		outputFile = filepath.Join(path, file)
	}

	if path != "" {
		err = os.MkdirAll(path, 0o755)
		if err != nil {
			return fmt.Errorf("oops! error creating path\n%v", err)
		}
	}
	return nil
}
