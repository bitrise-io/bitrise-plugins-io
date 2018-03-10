package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/spf13/cobra"
)

var setAuthTokenCmd = &cobra.Command{
	Use:   "set-auth-token",
	Short: "Set API authentication token",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := setAuthToken(); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setAuthTokenCmd)
}

func setAuthToken() error {
	log.Infof("")
	log.Infof("Set authentication token...")

	if len(os.Args) != 1 {
		return errors.New("Failed to set authentication token, error: invalid number of arguments")
	}

	if err := configs.SetAPIToken(os.Args[0]); err != nil {
		return fmt.Errorf("Failed to set authentication token, error: %s", err)
	}

	log.Infof("Authentication token set successfully...")

	err := services.ValidateAuthToken()
	if err != nil {
		return err
	}
	log.Infof("Authentication token validated successfully...")
	return nil
}
