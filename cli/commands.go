package cli

import (
	"errors"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/slapec93/bitrise-plugins-io/configs"
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

func setAuthToken(c *cli.Context) {
	log.Infof("")
	log.Infof("\x1b[34;1mSet authentication token...\x1b[0m")

	args := c.Args()
	if len(args) != 1 {
		log.Errorf("Failed to set authentication token, error: %s", errors.New("invalid number of arguments"))
		os.Exit(1)
	}

	if err := configs.SetAPIToken(args[0]); err != nil {
		log.Errorf("Failed to set authentication token, error: %s", err)
		os.Exit(1)
	}

	err := services.ValidateAuthToken()
	if err != nil {
		log.Errorf("\x1b[33;1m%s\x1b[0m", err)
	} else {
		log.Infof("\x1b[32;1mAuthentication token set successfully...\x1b[0m")
	}
}

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
