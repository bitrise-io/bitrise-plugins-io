package cmd

import (
	"fmt"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
)

func openURL(urlToOpen string) {
	fmt.Println(colorstring.Yellow("Opening URL:"), urlToOpen)
	if err := command.NewWithStandardOuts("open", urlToOpen).Run(); err != nil {
		log.Printf("Failed to open URL in browser: %s", err)
	}
}
