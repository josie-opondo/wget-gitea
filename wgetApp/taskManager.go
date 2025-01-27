package wgetApp

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

	return nil
}
