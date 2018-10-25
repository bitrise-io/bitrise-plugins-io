package cmd

import (
	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// buildTriggerCmd represents the trigger command
var buildTriggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger build",
	Long:  `Trigger a new build`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(triggerBuild())
	},
}

func init() {
	appsCmd.AddCommand(buildTriggerCmd)
}

// BuildTriggerResponseModel ...
type BuildTriggerResponseModel struct {
	Status            string `json:"status"`
	Message           string `json:"message"`
	Slug              string `json:"slug"`
	Service           string `json:"service"`
	BuildSlug         string `json:"build_slug"`
	BuildNumber       int    `json:"build_number"`
	BuildURL          string `json:"build_url"`
	TriggeredWorkflow string `json:"triggered_workflow"`
}

// BuildTriggerParamsModel ...
type BuildTriggerParamsModel struct {
	HookInfo struct {
		Type              string `json:"type"`
		BuildTriggerToken string `json:"build_trigger_token"`
	} `json:"hook_info"`
	BuildParams struct {
		Branch       string `json:"branch"`
		WorkflowID   string `json:"workflow_id"`
		Environments []struct {
			MappedTo string `json:"mapped_to"`
			Value    string `json:"value"`
			IsExpand bool   `json:"is_expand"`
		} `json:"environments"`
	} `json:"build_params"`
	TriggeredBy string `json:"triggered_by"`
}

func triggerBuild() error {
	params := map[string]string{}

	response, err := services.GetBitriseArtifacts(artifactsAppIDFlag, artifactsBuildIDFlag, params)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return NewRequestFailedError(response)
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != formatJSON, &ArtifactsListReponseModel{}))
}
