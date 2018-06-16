package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	bitriseioBaseURL = "https://www.bitrise.io"
)

var (
	browseAppSlugFlag   string
	browseBuildSlugFlag string
	openInBrowser       bool
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse/open bitrise.io website at the app",
	RunE: func(cmd *cobra.Command, args []string) error {
		urlToOpen := bitriseioBaseURL
		if browseAppSlugFlag != "" {
			urlToOpen = fmt.Sprintf("%s/app/%s", bitriseioBaseURL, browseAppSlugFlag)
			if browseBuildSlugFlag != "" {
				// In the future the URL will include both the app & the build ID
				// but right now it's not required and not even an option.
				urlToOpen = fmt.Sprintf("%s/builds/%s", bitriseioBaseURL, browseBuildSlugFlag)
			}
		}
		openURL(urlToOpen, openInBrowser)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.Flags().StringVarP(&browseAppSlugFlag, "app", "a", "", "Slug of the app")
	browseCmd.Flags().StringVarP(&browseBuildSlugFlag, "build", "b", "", "Slug of the build")
	browseCmd.Flags().BoolVar(&openInBrowser, "open", true, "Open in browser? If set to false it'll only print out the URL but will not open it in browser.")
}
