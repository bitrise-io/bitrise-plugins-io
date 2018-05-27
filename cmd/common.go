package cmd

import (
	"fmt"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
)

func openURL(urlToOpen string, openInBrowser bool) {
	if openInBrowser {
		fmt.Println(colorstring.Yellow("Opening URL:"), urlToOpen)
		if err := command.NewWithStandardOuts("open", urlToOpen).Run(); err != nil {
			log.Errorf("Failed to open the specified URL in browser: %s", err)
		}
	} else {
		fmt.Println(colorstring.Yellow("URL:"), urlToOpen)
	}
}
