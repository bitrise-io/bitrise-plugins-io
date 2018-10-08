package cmd

import (
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "Logout",
	Long:    `Logout from plugin`,
	Example: `logout`,
	RunE:    unauth,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
