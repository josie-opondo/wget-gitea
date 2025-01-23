package wgetApp

import (
	"fmt"
	"os"
	"strings"
	wgetutils "wget/wgetUtils"
)

// parser processes command-line arguments and configures WgetApp settings.
func (app *WgetApp) parser() error {
	mirrorMode := false // Flag to track if --mirror is used
	track := false      // Flag to track if a source file is provided (-i=)

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-O=") {
			app.urlArgs.file = arg[len("-O="):]
		} else if strings.HasPrefix(arg, "-P=") {
			app.urlArgs.path = arg[len("-P="):]
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			if err := wgetutils.RateLimitValidator(arg); err != nil {
				return err
			}
			app.urlArgs.rateLimit = arg[len("--rate-limit="):]
		} else if strings.HasPrefix(arg, "--mirror") {
			app.urlArgs.mirroring = true
			mirrorMode = true
		} else if strings.HasPrefix(arg, "--convert-links") {
			if !mirrorMode {
				return fmt.Errorf("error: --convert-links can only be used with --mirror")
			}
			app.urlArgs.convertLinksFlag = true
		} else if strings.HasPrefix(arg, "-R=") || strings.HasPrefix(arg, "--reject=") {
			if !mirrorMode {
				return fmt.Errorf("error: --reject can only be used with --mirror")
			}
			if strings.HasPrefix(arg, "-R=") {
				app.urlArgs.rejectFlag = arg[len("-R="):]
			} else {
				app.urlArgs.rejectFlag = arg[len("--reject="):]
			}
		} else if strings.HasPrefix(arg, "-X=") || strings.HasPrefix(arg, "--exclude=") {
			if !mirrorMode {
				return fmt.Errorf("error: --exclude can only be used with --mirror")
			}
			if strings.HasPrefix(arg, "-X=") {
				app.urlArgs.excludeFlag = arg[len("-X="):]
			} else {
				app.urlArgs.excludeFlag = arg[len("--exclude="):]
			}
		} else if strings.HasPrefix(arg, "-B") {
			app.urlArgs.workInBackground = true
		} else if strings.HasPrefix(arg, "-i=") {
			app.urlArgs.sourceFile = arg[len("-i="):]
			track = true
		} else if strings.HasPrefix(arg, "http") {
			app.urlArgs.url = arg
		} else {
			return fmt.Errorf("error: Unrecognized argument '%s'", arg)
		}
	}

	// Validate rate limit format (must end with 'k' or 'm')
	if app.urlArgs.rateLimit != "" {
		lastChar := strings.ToLower(string(app.urlArgs.rateLimit[len(app.urlArgs.rateLimit)-1]))
		if lastChar != "k" && lastChar != "m" {
			return fmt.Errorf("invalid rateLimit")
		}
	}

	// Validate background mode restrictions
	if app.urlArgs.workInBackground {
		if app.urlArgs.sourceFile != "" || app.urlArgs.path != "" {
			return fmt.Errorf("-B flag should not be used with -i or -P flags")
		}
	}

	// Ensure --mirror is not combined with incompatible flags
	if app.urlArgs.mirroring {
		if app.urlArgs.file != "" || app.urlArgs.path != "" || app.urlArgs.rateLimit != "" ||
			app.urlArgs.sourceFile != "" || app.urlArgs.workInBackground {
			return fmt.Errorf("error: --mirror can only be used with --convert-links, --reject, --exclude, and a URL. No other flags are allowed")
		}
	} else {
		// Ensure --convert-links, --reject, and --exclude are only used with --mirror
		if app.urlArgs.convertLinksFlag || app.urlArgs.rejectFlag != "" || app.urlArgs.excludeFlag != "" {
			return fmt.Errorf("error: --convert-links, --reject, and --exclude can only be used with --mirror")
		}
	}

	// Ensure a URL or source file is provided for valid execution
	if app.urlArgs.url == "" && !track {
		return fmt.Errorf("error: URL not provided")
	}

	// Validate the url
	err := wgetutils.ValidateURL(app.urlArgs.url)
	if err != nil {
		return fmt.Errorf("error: invalid url provided")
	}

	return nil
}
