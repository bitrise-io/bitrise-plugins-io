package services

import "fmt"

// BuildsReponseModel ...
type BuildsReponseModel struct {
	Slug                    string `json:"slug"`
	Status                  int    `json:"status"`
	StatusText              string `json:"status_text"`
	IsOnHold                bool   `json:"is_on_hold"`
	BuildNumber             int    `json:"build_number"`
	Branch                  string `json:"branch"`
	Tag                     string `json:"tag"`
	PullRequestTargetBranch string `json:"pull_request_target_branch"`
	PullRequestID           int    `json:"pull_request_id"`
	CommitHash              string `json:"commit_hash"`
	CommitMessage           string `json:"commit_message"`
	TriggeredWorkflow       string `json:"triggered_workflow"`
	TriggeredBy             string `json:"triggered_by"`
}

// BuildsListReponseModel ...
type BuildsListReponseModel struct {
	Data []BuildsReponseModel `json:"data"`
}

// TriggerInfoString ...
func (respModel *BuildsReponseModel) TriggerInfoString() string {
	if respModel.PullRequestID > 0 {
		return fmt.Sprintf("(#%d) %s > %s", respModel.PullRequestID, respModel.Branch, respModel.PullRequestTargetBranch)
	}
	if len(respModel.Tag) > 0 {
		return fmt.Sprintf("tag: %s", respModel.Tag)
	}
	return fmt.Sprintf("push: %s", respModel.Branch)
}
