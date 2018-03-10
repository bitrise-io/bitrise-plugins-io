package cli

import (
	"os"

	"github.com/urfave/cli"
)

var commands = []cli.Command{
	setAuthTokenCmd,
	appsCmd,
	buildsCmd,
}

var (
	nextFlag = cli.StringFlag{
		Name:   "next",
		Usage:  "Next parameter for paging",
		EnvVar: "NEXT",
	}
	limitFlag = cli.Int64Flag{
		Name:   "limit",
		Usage:  "Limit parameter for paging",
		EnvVar: "LIMIT",
	}
	sortByFlag = cli.StringFlag{
		Name:   "sort_by",
		Usage:  "Sort by parameter for listing",
		EnvVar: "SORT_BY",
	}
)

func getFlag(c *cli.Context, envName, flagName string) string {
	flagValue := c.String(flagName)
	if len(flagValue) == 0 {
		flagValue = os.Getenv(envName)
	}
	return flagValue
}

func fetchFlagsForObjectListing(c *cli.Context) map[string]string {
	return map[string]string{
		"next":    getFlag(c, "NEXT", "next"),
		"limit":   getFlag(c, "LIMIT", "limit"),
		"sort_by": getFlag(c, "SORT_BY", "sort_by"),
	}
}
