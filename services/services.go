package services

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
)

const (
	apiRootURL = "https://api.bitrise.io/v0.1"
)

func getBitriseObjectList(subURL string, params map[string]string) error {
	req, err := getRequest(fmt.Sprintf("%s/%s", apiRootURL, subURL), params)
	if err != nil {
		return errors.WithStack(err)
	}

	client := createClient()
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
		return fmt.Errorf("fetching apps from Bitrise IO, failed with status code: %d: %s", resp.StatusCode, resp.Status)
	}

	response := map[string]interface{}{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.WithStack(err)
	}

	prettyResp, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}
	log.Infof(string(prettyResp))

	return nil
}

// GetBitriseAppsForUser ...
func GetBitriseAppsForUser(params map[string]string) error {
	return getBitriseObjectList("apps", params)
}

// GetBitriseBuildsForApp ...
func GetBitriseBuildsForApp(appSlug string, params map[string]string) error {
	return getBitriseObjectList(fmt.Sprintf("apps/%s/builds", appSlug), params)
}

// ValidateAuthToken ...
func ValidateAuthToken() error {
	req, err := getRequest(fmt.Sprintf("%s/me", apiRootURL), map[string]string{})
	if err != nil {
		return errors.WithStack(err)
	}

	client := createClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request, error: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 210 {
		return fmt.Errorf("Invalid authentication token")
	}
	return nil
}
