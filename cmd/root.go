package cmd

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bitrise-plugins-io",
	Short: "Bitrise IO plugin",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(-1)
	}
}

func init() {
	configs.DataDir = os.Getenv("BITRISE_PLUGIN_INPUT_DATA_DIR")
}
