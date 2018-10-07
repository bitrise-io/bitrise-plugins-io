package cmd

import (
	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var unauthCmd = &cobra.Command{
	Use:     "unauth",
	Short:   "Unauthenticate",
	Long:    `Unauthenticate`,
	Example: `unauth`,
	RunE:    unauth,
}

func init() {
	rootCmd.AddCommand(unauthCmd)
	unauthCmd.Flags().StringVar(&formatFlag, "format", "pretty", "Output format, one of: [pretty, json]")
}

func unauth(cmd *cobra.Command, args []string) error {
	if err := configs.SetAPIToken(""); err != nil {
		return errors.Errorf("Failed to unauthenticate, error: %s", err)
	}

	return nil
}
