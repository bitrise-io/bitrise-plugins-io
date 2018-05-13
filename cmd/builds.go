package cmd

import (
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/spf13/cobra"
)

var (
	appSlugFlag string
)

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Get builds for app",
	Run: func(cmd *cobra.Command, args []string) {
		if err := builds(); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildsCmd)
	buildsCmd.Flags().StringVar(&nextFlag, "next", "", "Next parameter for paging")
	buildsCmd.Flags().StringVarP(&limitFlag, "limit", "l", "", "Limit parameter for paging")
	buildsCmd.Flags().StringVar(&sortFlag, "sort", "", "Sort by parameter for listing")
	buildsCmd.Flags().StringVarP(&appSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
	buildsCmd.Flags().StringVar(&formatFlag, "format", "pretty", "Output format, one of: [pretty, json]")
}

func builds() error {
	params := map[string]string{
		"next":    nextFlag,
		"limit":   limitFlag,
		"sort_by": sortFlag,
	}

	response, err := services.GetBitriseBuildsForApp(appSlugFlag, params)
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
