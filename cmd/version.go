package cmd

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVar(&formatFlag, "format", "pretty", "Output format, one of: [pretty, json]")
}

func printVersion() {
	if formatFlag == "json" {
		fmt.Printf(fmt.Sprintf(`{"data":"%s"}`, version.VERSION))
	} else {
		fmt.Println(version.VERSION)
	}
}
