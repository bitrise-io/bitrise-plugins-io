package cmd

import (
	"fmt"
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
	buildsCmd.Flags().StringVarP(&nextFlag, "next", "", "", "Next parameter for paging")
	buildsCmd.Flags().StringVarP(&limitFlag, "limit", "", "", "Limit parameter for paging")
	buildsCmd.Flags().StringVarP(&sortByFlag, "sort_by", "", "", "Sort by parameter for listing")
	buildsCmd.Flags().StringVarP(&appSlugFlag, "app_slug", "", "", "Slug of the app where the builds belong to")
}

func builds() error {
	params := map[string]string{
		"next":    nextFlag,
		"limit":   limitFlag,
		"sort_by": sortByFlag,
	}

	response, err := services.GetBitriseBuildsForApp(appSlugFlag, params)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}
