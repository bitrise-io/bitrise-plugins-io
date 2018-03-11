package cmd

import (
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	nextFlag   string
	limitFlag  string
	sortFlag   string
	formatFlag string
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Get apps for user",
	Run: func(cmd *cobra.Command, args []string) {
		if err := apps(); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.Flags().StringVar(&nextFlag, "next", "", "Next parameter for paging")
	appsCmd.Flags().StringVar(&limitFlag, "limit", "", "Limit parameter for paging")
	appsCmd.Flags().StringVar(&sortFlag, "sort", "", "Sort by parameter for listing")
	appsCmd.Flags().StringVar(&formatFlag, "format", "pretty", "Output format, one of: [pretty, json]")
}

func apps() error {
	params := map[string]string{
		"next":    nextFlag,
		"limit":   limitFlag,
		"sort_by": sortFlag,
	}

	response, err := services.GetBitriseAppsForUser(params)
	if err != nil {
		return errors.WithStack(err)
	}

	if response.Error != "" {
		printErrorOutput(response.Error, formatFlag != "json")
		os.Exit(1)
		return nil
	}
	printOutput(response.Data, formatFlag != "json")
	return nil
}
