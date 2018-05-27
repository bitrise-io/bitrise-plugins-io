package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	bitriseioBaseURL = "https://www.bitrise.io"
)

var (
	browseAppSlugFlag string
	openInBrowser     bool
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse/open bitrise.io website at the app",
	RunE: func(cmd *cobra.Command, args []string) error {
		urlToOpen := bitriseioBaseURL
		if browseAppSlugFlag != "" {
			urlToOpen = fmt.Sprintf("%s/app/%s", bitriseioBaseURL, browseAppSlugFlag)
		}
		openURL(urlToOpen, openInBrowser)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.Flags().StringVarP(&browseAppSlugFlag, "app", "a", "", "Slug of the app")
	browseCmd.Flags().BoolVar(&openInBrowser, "open", true, "Open in browser? If set to false it'll only print out the URL but will not open it in browser.")
}
