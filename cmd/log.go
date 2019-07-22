package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Get log of build",
	Long:  `Get log of build`,
	RunE:  getLog,
}

var (
	logAppSlugFlag   string
	logBuildSlugFlag string
)

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.
	logCmd.Flags().StringVarP(&logAppSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
	logCmd.Flags().StringVarP(&logBuildSlugFlag, "build", "b", "", "Slug of the build where the log belong to")
}

func getLog(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		if len(args) > 1 {
			return errors.Errorf("Too many arguments (%+v), only a single one (APP-SLUG/BUILD-SLUG) is allowed.", args)
		}
		splits := strings.Split(args[0], "/")
		if len(splits) != 2 {
			return errors.Errorf("Invalid argument (%+v), should be in format: APP-SLUG/BUILD-SLUG (e.g. 3...0/1...8)", splits)
		}
		logAppSlugFlag = splits[0]
		logBuildSlugFlag = splits[1]
	}
	fmt.Printf("App: %s | Build: %s\n", logAppSlugFlag, logBuildSlugFlag)

	fmt.Println("Retrieving Build and Log info ...")
	params := map[string]string{}
	response, err := services.GetBuildLogInfo(logAppSlugFlag, logBuildSlugFlag, params)
	if err != nil {
		return errors.WithStack(err)
	}

	if response.Error != "" {
		return services.NewRequestFailedError(response)
	}

	logInfo := struct {
		ExpiringRawLogURL string `json:"expiring_raw_log_url"`
		IsArchived        bool   `json:"is_archived"`
	}{}

	if err := json.Unmarshal(response.Data, &logInfo); err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("Downloading full log ...")
	fullLogData, err := services.LoadFullLog(logInfo.ExpiringRawLogURL)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("LOG: %s", fullLogData)
	return nil
}
