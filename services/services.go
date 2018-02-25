package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
	"github.com/slapec93/bitrise-plugins-io/configs"
)

const (
	apiRootURL = "https://api.bitrise.io/v0.1"
)

// GetBitriseAppsForUser ...
func GetBitriseAppsForUser() error {
	config, err := configs.ReadConfig()
	if err != nil {
		return errors.WithStack(err)
	}
	if len(config.BitriseAPIAuthenticationToken) < 1 {
		return errors.New("Bitrise API token isn't set, please set up with bitrise :io add-auth-token AUTH-TOKEN")
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/apps", apiRootURL), nil)
	if err != nil {
		return fmt.Errorf("failed to create request, error: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.BitriseAPIAuthenticationToken))

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request, error: %s", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body, error: %#v", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 210 {
		return fmt.Errorf("fetching apps from Bitrise IO, failed with status code: %d", resp.StatusCode)
	}

	response := map[string]interface{}{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.WithStack(err)
	}

	prettyResp, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		return errors.WithStack(err)
	}
	log.Infof(string(prettyResp))

	return nil
}
