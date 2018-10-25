package cmd

import (
	"github.com/spf13/cobra"
)

var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Browse/open Github issue tracker",
	Long:  "Browse/open Github issue tracker of source repository to check/report issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		openURL("https://github.com/bitrise-core/bitrise-plugins-io/issues", openInBrowser)
		return nil
	},
	Aliases: []string{"issues", "bug-report"},
}

func init() {
	rootCmd.AddCommand(feedbackCmd)
	feedbackCmd.Flags().BoolVar(&openInBrowser, "open", true, "Open in browser? If set to false it'll only print out the URL but will not open it in browser.")
}
