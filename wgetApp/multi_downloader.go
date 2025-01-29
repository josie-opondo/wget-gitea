package wgetApp


/*
downloadMultipleFiles 
- Opens the file containing the URLs.
- Reads URLs line by line, skipping empty lines.
- Uses a WaitGroup to manage multiple concurrent downloads.
-  Calls AsyncDownload to handle each URL download asynchronously.
- Waits for all goroutines to complete before returning.
*/
func (app *wgetApp) downloadMultipleFiles(filePath, outputFile, limit, directory string) error {
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
			err := app.AsyncDownload(outputFile, url, limit, directory)
			if err != nil {
				return err
			}
			return nil
		}(url)
	}
	wg.Wait()

	return nil
}
