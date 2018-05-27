package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	formatJSON   = "json"
	formatPretty = "pretty"
)

var (
	formatFlag string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bitrise-plugins-io",
	Short: "Command Line User Interface for bitrise.io",
	Long: `Command Line User Interface for bitrise.io

Uses the official Bitrise API (v0.1 docs: https://devcenter.bitrise.io/api/v0.1/ )
`,
	Example: `  If you use it as a Bitrise CLI plugin:
    $ bitrise :io [command]

  If you use it as a stand-alone tool:
    $ env BITRISE_PLUGIN_INPUT_DATA_DIR=/path/where/login/data/should/be/stored \
      bitrise-plugins-io [command]`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if inputErr, ok := errors.Cause(err).(*InputError); ok {
			// Input Error
			fmt.Printf(colorstring.Red("INPUT ERROR:")+" %s\n", inputErr)
		} else if confErr, ok := errors.Cause(err).(*services.ConfigError); ok {
			// Request Config Error (missing Personal Access Token)
			printErrorOutput(confErr.Error(), formatFlag == formatPretty)
		} else if reqFailErr, ok := errors.Cause(err).(*RequestFailedError); ok {
			// Request Failed (non successful response) Error
			response := reqFailErr.Response
			if response.StatusCode == http.StatusUnauthorized {
				if formatFlag == formatPretty {
					log.Warnf("Unauthorized - your Personal Access Token most likely expired or was revoked. Use the auth command to re-authenticate.")
				}
				if err := configs.SetAPIToken(""); err != nil {
					log.Errorf("Failed to clear stored Personal Access Token: %+v", err)
				}
			} else {
				printErrorOutput(response.Error, formatFlag == formatPretty)
			}
		} else {
			// Any other error: print with stack trace
			fmt.Printf("%+v\n", err)
		}
		os.Exit(-1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configs.DataDir, "bitrise-plugin-input-data-dir", "", "Bitrise Plugin data dir ($BITRISE_PLUGIN_INPUT_DATA_DIR)")
	configs.DataDir = envutil.GetenvWithDefault("BITRISE_PLUGIN_INPUT_DATA_DIR", configs.DataDir)

	rootCmd.PersistentFlags().StringVar(&configs.APIRootURL, "api-root-url", "https://api.bitrise.io/v0.1", "API root URL ($BITRISE_API_ROOT_URL)")
	configs.APIRootURL = envutil.GetenvWithDefault("BITRISE_API_ROOT_URL", configs.APIRootURL)

	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", formatPretty, "Output format, one of: [pretty, json]")
}
