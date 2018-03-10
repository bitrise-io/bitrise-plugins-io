package cli

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

var appsCmd = cli.Command{
	Name:   "apps",
	Usage:  "Get apps for user",
	Action: apps,
	Flags: []cli.Flag{
		nextFlag,
		limitFlag,
		sortByFlag,
	},
}

func apps(c *cli.Context) {
	params := fetchFlagsForObjectListing(c)
	response, err := services.GetBitriseAppsForUser(params)
	if err != nil {
		log.Errorf("Failed to fetch application list, error: %s", err)
		os.Exit(1)
	}
	fmt.Println(response)
}
