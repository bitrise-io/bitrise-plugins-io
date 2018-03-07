package cli

import "github.com/urfave/cli"

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
)
