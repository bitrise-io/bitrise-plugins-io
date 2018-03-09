package cli

import (
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/slapec93/bitrise-plugins-io/services"
	"github.com/urfave/cli"
)

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

//=======================================
// Actions
//=======================================

func apps(c *cli.Context) {
	err := services.GetBitriseAppsForUser(fetchFlagsForObjectListing(c))
	if err != nil {
		log.Errorf("Failed to fetch application list, error: %s", err)
		os.Exit(1)
	}
}

func builds(c *cli.Context) {
	appSlug := getFlag(c, "APP_SLUG", "app-slug")
	err := services.GetBitriseBuildsForApp(appSlug, fetchFlagsForObjectListing(c))
	if err != nil {
		log.Errorf("Failed to fetch build list, error: %s", err)
		os.Exit(1)
	}
}
