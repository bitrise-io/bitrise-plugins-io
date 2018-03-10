package cli

import (
	"errors"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

var setAuthTokenCmd = cli.Command{
	Name:   "set-auth-token",
	Usage:  "Set API authentication token",
	Action: setAuthToken,
}

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

	log.Infof("\x1b[32;1mAuthentication token set successfully...\x1b[0m")

	err := services.ValidateAuthToken()
	if err != nil {
		log.Errorf("\x1b[33;1m%s...\x1b[0m", err)
	} else {
		log.Infof("\x1b[32;1mAuthentication token validated successfully...\x1b[0m")
	}
}
