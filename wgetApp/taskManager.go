package wgetApp

import "strings"

func (app *WgetApp) taskManager(err error) error {
	if err != nil {
		return err
	}

	// Mirror website handling
	if app.urlArgs.mirroring {
		err := app.mirror(app.urlArgs.url, app.urlArgs.rejectFlag, app.urlArgs.excludeFlag, app.urlArgs.convertLinksFlag)
		if err != nil {
			return err
		}
		return nil
	}

	// If no file name is provided, derive it from the url
	if app.urlArgs.file == "" && app.urlArgs.url != "" {
		urlParts := strings.Split(app.urlArgs.url, "/")
		app.urlArgs.file = urlParts[len(urlParts)-1]
	}

	// Handle the work-in-background flag
	if app.urlArgs.workInBackground {
		err := app.downloadInBackground(app.urlArgs.file, app.urlArgs.url, app.urlArgs.rateLimit)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
