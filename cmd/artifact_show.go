package cmd

import (
	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// artifactShowCmd represents the show command
var artifactShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the specified build artifact",
	Long:  `Show the specified build artifact`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(artifactShow())
	},
}

var (
	artifactsShowArtifactIDFlag string
)

func init() {
	artifactsCmd.AddCommand(artifactShowCmd)

	artifactShowCmd.Flags().StringVar(&artifactsShowArtifactIDFlag, "slug", "", "Slug of the artifact to show")
}

func fetchArtifact(appID, buildID, artifactID string, params map[string]string) (services.Response, error) {
	response, err := services.GetBitriseArtifact(appID, buildID, artifactID, params)
	if err != nil {
		return services.Response{}, err
	}

	if response.Error != "" {
		return services.Response{}, NewRequestFailedError(response)
	}

	return response, nil
}

func artifactShow() error {
	params := map[string]string{}

	response, err := fetchArtifact(artifactsAppIDFlag, artifactsBuildIDFlag, artifactsShowArtifactIDFlag, params)
	if err != nil {
		return errors.WithStack(err)
	}

	printOutput(response.Data, formatFlag != formatJSON)
	return nil
}
