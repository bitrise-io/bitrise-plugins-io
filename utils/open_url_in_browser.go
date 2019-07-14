package utils

import (
	"fmt"
	"runtime"

	"github.com/bitrise-io/go-utils/command"
	"github.com/pkg/errors"
)

const (
	bitriseioBaseURL = "https://www.bitrise.io"
)

// OpenPageInBrowser ...
func OpenPageInBrowser(appID, buildID string) error {
	urlToOpen := GetURLForPage(appID, buildID)
	return OpenURLInBrowser(urlToOpen)
}

// GetURLForPage ...
func GetURLForPage(appID, buildID string) string {
	urlToOpen := bitriseioBaseURL
	if appID != "" {
		urlToOpen = fmt.Sprintf("%s/app/%s", bitriseioBaseURL, appID)
		if buildID != "" {
			// In the future the URL will include both the app & the build ID
			// but right now it's not required and not even an option.
			urlToOpen = fmt.Sprintf("%s/build/%s", bitriseioBaseURL, buildID)
		}
	}
	return urlToOpen
}

// OpenURLInBrowser ...
func OpenURLInBrowser(urlToOpen string) error {
	openCmd := "open"
	if runtime.GOOS == "linux" {
		openCmd = "xdg-open"
	}

	if outs, err := command.New(openCmd, urlToOpen).RunAndReturnTrimmedCombinedOutput(); err != nil {
		return errors.Wrapf(err, "Failed to open the specified URL in browser: %s", outs)
	}
	return nil
}
