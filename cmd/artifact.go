package cmd

import (
	"github.com/spf13/cobra"
)

// artifactsCmd represents the artifact command
var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "Manage artifact",
	Long:  `Manage artifact`,
}

var (
	artifactsAppIDFlag   string
	artifactsBuildIDFlag string
)

func init() {
	rootCmd.AddCommand(artifactsCmd)

	artifactsCmd.PersistentFlags().StringVarP(&artifactsAppIDFlag, "app", "a", "", "App ID (slug)")
	artifactsCmd.PersistentFlags().StringVarP(&artifactsBuildIDFlag, "build", "b", "", "Build ID (slug)")
}
