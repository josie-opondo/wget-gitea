package appstate

func (app *AppState) downloadInBackground(file, urlStr, rateLimit string) error{
	// Parse the URL to derive the output name
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL")
	}
	outputName := filepath.Base(parsedURL.Path) // Get the file name from the URL
	if file != "" {
		outputName = file
	}
	path := "." // Default path to save the file
	// Create the wget-log file to log output
	logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("error creating log file:\n%v", err)
	}
	defer logFile.Close()

	// Ensure the output directory exists
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory:\n%v", err)
	}
	cmd := exec.Command(os.Args[0], "-O="+outputName, "-P="+path, "--rate-limit="+rateLimit, urlStr)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	fmt.Println("Output will be written to \"wget-log\".")

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting download:\n%v", err)

	}



}