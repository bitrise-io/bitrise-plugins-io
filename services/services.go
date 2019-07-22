package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bitrise-io/bitrise-plugins-io/configs"
	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
)

// Response ...
type Response struct {
	Data       []byte
	Error      string
	StatusCode int
}

func newErrorResponse(statusCode int, errMsg string) Response {
	return Response{StatusCode: statusCode, Error: errMsg}
}

func newSuccessResponse(statusCode int, bodyData []byte) Response {
	return Response{Data: bodyData}
}

func wrapResponse(response *http.Response) (Response, error) {
	if response.StatusCode < 200 || response.StatusCode > 210 {
		body := map[string]interface{}{}
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			return Response{}, errors.WithStack(err)
		}

		return newErrorResponse(response.StatusCode, fmt.Sprintf("%s", body["message"])), nil
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, errors.WithStack(err)
	}

	return newSuccessResponse(response.StatusCode, data), nil
}

func bitriseGetRequest(subURL string, params map[string]string) (Response, error) {
	req, err := request("GET", fmt.Sprintf("%s/%s", configs.APIRootURL, subURL), params, nil)
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

func bitrisePostRequest(subURL string, requestBody map[string]interface{}) (Response, error) {
	req, err := request("POST", fmt.Sprintf("%s/%s", configs.APIRootURL, subURL), nil, requestBody)
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

// AppSortBy ...
type AppSortBy string

const (
	// SortAppsByCreatedAt ...
	SortAppsByCreatedAt AppSortBy = "created_at"
	// SortAppsByLastBuildAt ...
	SortAppsByLastBuildAt AppSortBy = "last_build_at"
)

// GetBitriseAppsForUser ...
func GetBitriseAppsForUser(next, limit string, sortBy AppSortBy, title string) (Response, error) {
	return bitriseGetRequest("apps", map[string]string{
		"next":    next,
		"limit":   limit,
		"sort_by": string(sortBy),
		"title":   title,
	})
}

// GetBitriseBuildsForApp ...
func GetBitriseBuildsForApp(appSlug string, params map[string]string) (Response, error) {
	return bitriseGetRequest(fmt.Sprintf("apps/%s/builds", appSlug), params)
}

// GetBitriseArtifacts ...
func GetBitriseArtifacts(appSlug string, buildSlug string, params map[string]string) (Response, error) {
	return bitriseGetRequest(fmt.Sprintf("apps/%s/builds/%s/artifacts", appSlug, buildSlug), params)
}

// GetBitriseArtifact ...
func GetBitriseArtifact(appSlug, buildSlug, artifactSlug string, params map[string]string) (Response, error) {
	return bitriseGetRequest(fmt.Sprintf("apps/%s/builds/%s/artifacts/%s", appSlug, buildSlug, artifactSlug), params)
}

// GetBuildLogInfo ...
func GetBuildLogInfo(appSlug, buildSlug string, params map[string]string) (Response, error) {
	return bitriseGetRequest(fmt.Sprintf("apps/%s/builds/%s/log", appSlug, buildSlug), params)
}

// LoadFullLog ...
func LoadFullLog(fullLogURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", fullLogURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send request")
	}
	defer responseBodyCloser(resp)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// responseBodyCloser closes a HTTP response body with logging the error
func responseBodyCloser(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		// TODO: modify this to use logrus
		log.Printf("Failed to close response body: %+v", errors.WithStack(err))
	}
}

// ValidateAuthToken ...
func ValidateAuthToken() (Response, error) {
	return bitriseGetRequest("me", nil)
}

// RegisterRepository ...
func RegisterRepository(repoURL string) (Response, error) {
	return bitrisePostRequest("apps/register", map[string]interface{}{"repo_url": repoURL})
}

// RegisterSSHKey ...
func RegisterSSHKey(appSlug, publicKey, privateKey string) (Response, error) {
	params := map[string]interface{}{
		"auth_ssh_private_key": privateKey,
		"auth_ssh_public_key":  publicKey,
	}
	return bitrisePostRequest(fmt.Sprintf("apps/%s/register-ssh-key", appSlug), params)
}

// RegisterWebhook ...
func RegisterWebhook(appSlug string) (Response, error) {
	return bitrisePostRequest(fmt.Sprintf("apps/%s/register-webhook", appSlug), nil)
}

// FinishAppRegistration ...
func FinishAppRegistration(appSlug, projectType, stackID string, organizationSlug *string, envs map[string]string, config map[string]interface{}) (Response, error) {
	params := map[string]interface{}{
		"mode":         "manual",
		"project_type": projectType,
		"stack_id":     stackID,
		"envs":         envs,
		"config":       config,
	}
	if organizationSlug != nil {
		params["oragenization_slug"] = *organizationSlug
	}
	return bitrisePostRequest(fmt.Sprintf("apps/%s/finish", appSlug), params)
}
