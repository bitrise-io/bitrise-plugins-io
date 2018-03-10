package cmd

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/log"
	"github.com/spf13/cobra"
)

var (
	nextFlag   string
	limitFlag  string
	sortByFlag string
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
	appsCmd.Flags().StringVarP(&nextFlag, "next", "", "", "Next parameter for paging")
	appsCmd.Flags().StringVarP(&limitFlag, "limit", "", "", "Limit parameter for paging")
	appsCmd.Flags().StringVarP(&sortByFlag, "sort_by", "", "", "Sort by parameter for listing")
}

func apps() error {
	params := map[string]string{
		"next":    nextFlag,
		"limit":   limitFlag,
		"sort_by": sortByFlag,
	}

	response, err := services.GetBitriseAppsForUser(params)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}
