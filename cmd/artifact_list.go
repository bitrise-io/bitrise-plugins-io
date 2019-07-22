package cmd

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	artifactsNextFlag  string
	artifactsLimitFlag string
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
	artifactListCmd.Flags().StringVarP(&artifactsNextFlag, "next", "n", "", "Next parameter for paging")
	artifactListCmd.Flags().StringVarP(&artifactsLimitFlag, "limit", "l", "", "Limit parameter for paging")
}

// ArtifactsListReponseModel ...
type ArtifactsListReponseModel struct {
	Data []struct {
		Title               string `json:"title"`
		ArtifactType        string `json:"artifact_type"`
		IsPublicPageEnabled bool   `json:"is_public_page_enabled"`
		Slug                string `json:"slug"`
		FileSizeBytes       int    `json:"file_size_bytes"`
	} `json:"data"`
}

// Pretty ...
func (respModel *ArtifactsListReponseModel) Pretty() string {
	linesOfTable := [][]string{}
	// headers
	linesOfTable = append(linesOfTable, []string{"Title", "Slug", "Size", "Pub Page Enabled?"})
	// data
	for _, aArtifact := range respModel.Data {
		linesOfTable = append(linesOfTable, []string{
			aArtifact.Title,
			aArtifact.Slug,
			fmt.Sprintf("%.2f KB", float64(aArtifact.FileSizeBytes)/1024),
			fmt.Sprintf("%t", aArtifact.IsPublicPageEnabled),
		})
	}

	return tabbedTableString(linesOfTable)
}

func artifactList() error {
	params := map[string]string{
		"next":  artifactsNextFlag,
		"limit": artifactsLimitFlag,
	}

	response, err := services.GetBitriseArtifacts(artifactsAppIDFlag, artifactsBuildIDFlag, params)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return services.NewRequestFailedError(response)
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != formatJSON, &ArtifactsListReponseModel{}))
}
