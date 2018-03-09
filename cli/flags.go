package cli

import (
	"os"

	"github.com/urfave/cli"
)

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
