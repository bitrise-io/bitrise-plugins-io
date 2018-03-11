package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-team/den/utils"
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
	buildSlugFlag string
)

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.
	logCmd.Flags().StringVarP(&appSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
	logCmd.Flags().StringVarP(&buildSlugFlag, "build", "b", "", "Slug of the build where the log belong to")
}

func loadFullLog(fullLogURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", fullLogURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send request")
	}
	defer utils.ResponseBodyCloser(resp)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
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
		appSlugFlag = splits[0]
		buildSlugFlag = splits[1]
	}
	fmt.Printf("App: %s | Build: %s\n", appSlugFlag, buildSlugFlag)

	serverURL := "https://api.bitrise.io/v0.1"

	config, err := configs.ReadConfig()
	if err != nil {
		return errors.WithStack(err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/apps/%s/builds/%s/log", serverURL, appSlugFlag, buildSlugFlag), nil)
	if err != nil {
		return errors.Wrap(err, "Failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.BitriseAPIAuthenticationToken))

	fmt.Println("Retrieving Build and Log info ...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to send request")
	}
	defer utils.ResponseBodyCloser(resp)

	if resp.StatusCode == 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.WithStack(err)
		}

		logInfo := struct {
			ExpiringRawLogURL string `json:"expiring_raw_log_url"`
			IsArchived        bool   `json:"is_archived"`
		}{}

		if err := json.Unmarshal(data, &logInfo); err != nil {
			return errors.WithStack(err)
		}

		fmt.Println("Downloading full log ...")
		fullLogData, err := loadFullLog(logInfo.ExpiringRawLogURL)
		if err != nil {
			return errors.WithStack(err)
		}

		fmt.Printf("LOG: %s", fullLogData)
	} else {
		log.Printf(colorstring.Red("RESPONSE:")+" %+v", resp)
	}
	fmt.Println()
	return nil
}
