package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/slapec93/bitrise-plugins-io/configs"
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
			} else {
				urlWithParams = fmt.Sprintf("%s&%s=%s", urlWithParams, paramName, paramValue)
			}
		}
	}
	return
}

// getRequest ...
func getRequest(url string, queryParams map[string]string) (*http.Request, error) {
	config, err := configs.ReadConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(config.BitriseAPIAuthenticationToken) < 1 {
		return nil, errors.New("Bitrise API token isn't set, please set up with bitrise :io add-auth-token AUTH-TOKEN")
	}
	req, err := http.NewRequest("GET", urlWithParameters(url, queryParams), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, error: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.BitriseAPIAuthenticationToken))

	return req, nil
}
