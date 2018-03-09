package cli

import "github.com/urfave/cli"

var commands = []cli.Command{
	cli.Command{
		Name:   "set-auth-token",
		Usage:  "Set API authentication token",
		Action: setAuthToken,
	},
	cli.Command{
		Name:   "apps",
		Usage:  "Get apps for user",
		Action: apps,
		Flags: []cli.Flag{
			nextFlag,
			limitFlag,
			sortByFlag,
		},
	},
	cli.Command{
		Name:   "builds",
		Usage:  "Get builds for app",
		Action: builds,
		Flags: []cli.Flag{
			nextFlag,
			limitFlag,
			sortByFlag,
			cli.StringFlag{
				Name:   "app-slug",
				Usage:  "Slug of the app where the builds belong to",
				EnvVar: "APP_SLUG",
			},
		},
	},
}
