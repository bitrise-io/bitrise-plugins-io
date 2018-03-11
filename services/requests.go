package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bitrise-core/bitrise-plugins-io/configs"
	"github.com/pkg/errors"
)

func createClient() http.Client {
	timeout := time.Duration(10 * time.Second)
	return http.Client{
		Timeout: timeout,
	}
}

func urlWithParameters(url string, queryParams map[string]string) (urlWithParams string) {
	isFirstParam := true
	urlWithParams = url
	for paramName, paramValue := range queryParams {
		if len(paramValue) > 0 {
			if isFirstParam {
				urlWithParams = fmt.Sprintf("%s?%s=%s", urlWithParams, paramName, paramValue)
				isFirstParam = false
			} else {
				urlWithParams = fmt.Sprintf("%s&%s=%s", urlWithParams, paramName, paramValue)
			}
		}
	}
	return
}

func request(method, url string, queryParams map[string]string, requestBody map[string]interface{}) (*http.Request, error) {
	config, err := configs.ReadConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(config.BitriseAPIAuthenticationToken) < 1 {
		return nil, errors.New("Bitrise API token isn't set, please set up with bitrise :io auth AUTH-TOKEN")
	}
	requestBytes, err := json.Marshal(requestBody)
	req, err := http.NewRequest(method, urlWithParameters(url, queryParams), bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, errors.Errorf("failed to create request, error: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.BitriseAPIAuthenticationToken))

	return req, nil
}
