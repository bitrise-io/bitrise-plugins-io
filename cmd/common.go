package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

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
