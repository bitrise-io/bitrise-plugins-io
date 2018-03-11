package cmd

import (
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	flagAPIToken string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Errorf("More than one argument specified (%+v), only one (the API Token) should be", args)
			os.Exit(1)
		}

		tokenToAuth := ""
		if flagAPIToken != "" {
			if len(args) > 0 {
				log.Errorf("Both the --token flag as well as a token arg (%+v) is specified. Only one should be", args)
				os.Exit(1)
			}
			tokenToAuth = flagAPIToken
		} else if len(args) == 1 {
			tokenToAuth = args[0]
		}

		if err := auth(tokenToAuth); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVar(&flagAPIToken, "token", "", "Authentication token")
	authCmd.Flags().StringVar(&formatFlag, "format", "pretty", "Output format, one of: [pretty, json]")
}

func auth(apiToken string) error {
	if apiToken == "" {
		return errors.New("Failed to set authentication token, error: no API Token specified")
	}

	if err := configs.SetAPIToken(apiToken); err != nil {
		return errors.Errorf("Failed to set authentication token, error: %s", err)
	}

	response, err := services.ValidateAuthToken()
	if err != nil {
		return err
	}

	if response.Error != "" {
		printErrorOutput(response.Error, formatFlag != "json")
		os.Exit(1)
		return nil
	}
	printOutput(response.Data, formatFlag != "json")
	return nil
}
