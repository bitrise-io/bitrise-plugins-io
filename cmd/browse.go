package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	bitriseioBaseURL = "https://www.bitrise.io"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse/open bitrise.io website at the app",
	RunE: func(cmd *cobra.Command, args []string) error {
		urlToOpen := bitriseioBaseURL
		if appSlugFlag != "" {
			urlToOpen = fmt.Sprintf("%s/app/%s", bitriseioBaseURL, appSlugFlag)
		}
		openURL(urlToOpen)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.Flags().StringVarP(&appSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
}
