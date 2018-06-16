package cmd

import (
	"github.com/spf13/cobra"
)

// artifactCmd represents the artifact command
var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage artifact",
	Long:  `Manage artifact`,
}

var (
	artifactAppIDFlag   string
	artifactBuildIDFlag string
)

func init() {
	rootCmd.AddCommand(artifactCmd)

	artifactCmd.PersistentFlags().StringVarP(&artifactAppIDFlag, "app", "a", "", "App ID (slug)")
	artifactCmd.PersistentFlags().StringVarP(&artifactBuildIDFlag, "build", "b", "", "Build ID (slug)")
}
