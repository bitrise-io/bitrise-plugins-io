package cmd

import (
	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	personalAccessTokenRegURL = "https://www.bitrise.io/me/profile#/security"
)

var (
	flagAPIToken string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the specified --token",
	Long: `Authenticate with the specified --token

` + colorstring.Green("You can register a new Personal Access Token at:") + " " + personalAccessTokenRegURL,
	Example: `auth --token=B1triseI0PersonalAccessT0ken`,
	RunE:    auth,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVar(&flagAPIToken, "token", "", "Personal Access token")
	authCmd.Flags().StringVar(&formatFlag, "format", "pretty", "Output format, one of: [pretty, json]")
}

func auth(cmd *cobra.Command, args []string) error {
	if flagAPIToken == "" {
		return NewInputError("No Personal Access Token specified. Register one at: " + personalAccessTokenRegURL)
	}

	if err := configs.SetAPIToken(flagAPIToken); err != nil {
		return errors.Errorf("Failed to set authentication token, error: %s", err)
	}

	return errors.WithStack(whoami())
}
