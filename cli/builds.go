package cli

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

var buildsCmd = cli.Command{
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
}

func builds(c *cli.Context) {
	appSlug := getFlag(c, "APP_SLUG", "app-slug")
	response, err := services.GetBitriseBuildsForApp(appSlug, fetchFlagsForObjectListing(c))
	if err != nil {
		log.Errorf("Failed to fetch build list, error: %s", err)
		os.Exit(1)
	}
	fmt.Println(response)
}
