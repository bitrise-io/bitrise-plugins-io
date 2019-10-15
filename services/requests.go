package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bitrise-io/bitrise-plugins-io/configs"
	"github.com/pkg/errors"
)

// ConfigError ...
type ConfigError struct {
	Err string
}

func (e *ConfigError) Error() string {
	return e.Err
}

// NewConfigError ...
func NewConfigError(err string) error {
	return &ConfigError{
		Err: err,
	}
}

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
		return nil, NewConfigError("Bitrise API token isn't set, please set it up with: $ bitrise :io auth --token=AUTH-TOKEN")
	}

	var bodyReader io.Reader
	if requestBody != nil {
		requestBytes, err := json.Marshal(requestBody)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		bodyReader = bytes.NewBuffer(requestBytes)
	}
	req, err := http.NewRequest(method, urlWithParameters(url, queryParams), bodyReader)
	if err != nil {
		return nil, errors.Errorf("failed to create request, error: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.BitriseAPIAuthenticationToken))

	return req, nil
}
