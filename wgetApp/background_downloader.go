package appstate

func (app *AppState) downloadInBackground(file, urlStr, rateLimit string) error{
	// Parse the URL to derive the output name
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL")
	}


}