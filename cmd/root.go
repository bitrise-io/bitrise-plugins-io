package cmd

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bitrise-plugins-io",
	Short: "Command Line User Interface for bitrise.io",
	Long: `Command Line User Interface for bitrise.io

--------------------------------------------------
If you use it as a Bitrise CLI plugin:
  $ bitrise :io [command]

If you use it as a stand-alone tool:
  $ env BITRISE_PLUGIN_INPUT_DATA_DIR=/path/where/login/data/should/be/stored \
    bitrise-plugins-io [command]
--------------------------------------------------
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if inputErr, ok := errors.Cause(err).(*InputError); ok {
			fmt.Printf(colorstring.Red("INPUT ERROR:")+" %s\n", inputErr)
		} else {
			// print with stack trace
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
}
