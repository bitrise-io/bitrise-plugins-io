package cmd

import (
	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// artifactListCmd represents the list command
var artifactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List artifacts",
	Long:  `List artifacts`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(artifactList())
	},
}

func init() {
	artifactsCmd.AddCommand(artifactListCmd)
}

func artifactList() error {
	params := map[string]string{}

	response, err := services.GetBitriseArtifacts(artifactsAppIDFlag, artifactsBuildIDFlag, params)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return NewRequestFailedError(response)
	}

	// return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &BuildsListReponseModel{}))
	printOutput(response.Data, formatFlag != formatJSON)
	return nil
}
