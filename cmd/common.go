package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

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

func tabbedTableString(linesOfFields [][]string) string {
	buf := bytes.NewBuffer([]byte{})
	prettyTabWriter := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)

	for _, aFieldsOfLine := range linesOfFields {
		if _, err := fmt.Fprintln(prettyTabWriter, strings.Join(aFieldsOfLine, "\t")); err != nil {
			panic(err)
		}
	}

	if err := prettyTabWriter.Flush(); err != nil {
		panic(err)
	}
	return buf.String()
}
