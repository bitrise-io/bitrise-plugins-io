package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
)

const (
	apiRootURL = "https://api.bitrise.io/v0.1"
)

// Response ...
type Response struct {
	Data  []byte
	Error string
}

func wrapResponse(response *http.Response) (Response, error) {
	if response.StatusCode < 200 || response.StatusCode > 210 {
		body := map[string]interface{}{}
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			return Response{}, errors.WithStack(err)
		}

		return Response{Error: fmt.Sprintf("%s", body["message"])}, nil
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{Data: data}, nil
}

func bitriseGetRequest(subURL string, params map[string]string) (Response, error) {
	req, err := getRequest(fmt.Sprintf("%s/%s", apiRootURL, subURL), params)
	if err != nil {
		return Response{}, errors.WithStack(err)
	}

	client := createClient()
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, errors.Errorf("failed to perform request, error: %s", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body, error: %#v", err)
		}
	}()

	return wrapResponse(resp)
}

// GetBitriseAppsForUser ...
func GetBitriseAppsForUser(params map[string]string) (Response, error) {
	return bitriseGetRequest("apps", params)
}

// GetBitriseBuildsForApp ...
func GetBitriseBuildsForApp(appSlug string, params map[string]string) (Response, error) {
	return bitriseGetRequest(fmt.Sprintf("apps/%s/builds", appSlug), params)
}

// ValidateAuthToken ...
func ValidateAuthToken() (Response, error) {
	return bitriseGetRequest("me", nil)
}
