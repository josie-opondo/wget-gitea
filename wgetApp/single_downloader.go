package wgetApp

import (
	"fmt"
	"os"
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
	// toDisplay, err := wgetutils.LoadShowProgressState(app.tempConfigFile)
	if err != nil {
		return err
	}
	fmt.Printf("started at %s\n", startTime.Format("2006-01-02 15:04:05"))

	resp, err := wgetutils.HttpRequest(fileURL)
	if err != nil {
		return fmt.Errorf("error downloading file:\nserver misbehaving")
	}
	defer resp.Body.Close()

	if path != "" {
		err = os.MkdirAll(path, 0o755)
		if err != nil {
			return fmt.Errorf("oops! error creating path\n%v", err)
		}
	}
	return nil
}
