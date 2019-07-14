package cmd

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/utils"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	browseAppSlugFlag       string
	browseBuildSlugFlag     string
	browseOpenInBrowserFlag bool
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse/open bitrise.io website at the app",
	RunE: func(cmd *cobra.Command, args []string) error {
		urlToOpen := utils.GetURLForPage(browseAppSlugFlag, browseBuildSlugFlag)
		if browseOpenInBrowserFlag {
			fmt.Println(colorstring.Yellow("Opening URL:"), urlToOpen)
			if err := utils.OpenURLInBrowser(urlToOpen); err != nil {
				return errors.WithStack(err)
			}
		} else {
			fmt.Println(colorstring.Yellow("URL:"), urlToOpen)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.Flags().StringVarP(&browseAppSlugFlag, "app", "a", "", "Slug of the app")
	browseCmd.Flags().StringVarP(&browseBuildSlugFlag, "build", "b", "", "Slug of the build")
	browseCmd.Flags().BoolVar(&browseOpenInBrowserFlag, "open", true, "Open in browser? If set to false it'll only print out the URL but will not open it in browser.")
}
