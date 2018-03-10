package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-core/bitrise-plugins-io/version"
	"github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

// Run ...
func Run() {
	cli.VersionPrinter = func(c *cli.Context) { fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version) }

	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "Bitrise IO plugin"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Before = func(c *cli.Context) error {
		configs.DataDir = os.Getenv("BITRISE_PLUGIN_INPUT_DATA_DIR")
		return nil
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Errorf("Finished with Error: %s", err)
		os.Exit(1)
	}
}
