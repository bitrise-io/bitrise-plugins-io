package formatter

import (
	"fmt"
	"strings"
)

// PrettyOneLinerText ...
func PrettyOneLinerText(msg string) string {
	s := fmt.Sprintf("%.50s", msg)         // print the first X chars
	s = strings.Replace(s, "\n", "↲", -1)  // replace newlines
	s = strings.Replace(s, "\r", "↲", -1)  // replace newlines
	s = strings.Replace(s, "\t", "  ", -1) // replace tabs
	return s
}
