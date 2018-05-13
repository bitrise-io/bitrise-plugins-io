package cmd

import (
	"fmt"

	"github.com/bitrise-io/go-utils/command"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse/open bitrise.io website at the app",
	RunE: func(cmd *cobra.Command, args []string) error {
		if appSlugFlag == "" {
			return errors.WithStack(command.NewWithStandardOuts("open", "https://www.bitrise.io").Run())
		}
		return errors.WithStack(command.NewWithStandardOuts("open", fmt.Sprintf("https://www.bitrise.io/app/%s", appSlugFlag)).Run())
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.Flags().StringVarP(&appSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
}
