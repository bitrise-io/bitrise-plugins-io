package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// artifactDownloadCmd represents the download command
var artifactDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download artifact",
	Long:  `Download artifact`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(artifactDownload())
	},
}

var (
	artifactsDownloadArtifactIDFlag string
)

func init() {
	artifactsCmd.AddCommand(artifactDownloadCmd)
	artifactDownloadCmd.Flags().StringVar(&artifactsDownloadArtifactIDFlag, "slug", "", "Slug of the artifact to download")
}

// ArtifactInfoResponseModel ...
type ArtifactInfoResponseModel struct {
	Data struct {
		ArtifactType         string `json:"artifact_type"`
		ExpiringDownloadURL  string `json:"expiring_download_url"`
		FileSizeBytes        int    `json:"file_size_bytes"`
		IsPublicPageEnabled  bool   `json:"is_public_page_enabled"`
		PublicInstallPageURL string `json:"public_install_page_url"`
		Slug                 string `json:"slug"`
		Title                string `json:"title"`
	} `json:"data"`
}

func artifactDownload() error {
	if artifactsDownloadArtifactIDFlag == "" {
		return NewInputError("No artifact ID specified.")
	}

	params := map[string]string{}
	response, err := fetchArtifact(artifactsAppIDFlag, artifactsBuildIDFlag, artifactsDownloadArtifactIDFlag, params)
	if err != nil {
		return errors.WithStack(err)
	}

	var infoModel ArtifactInfoResponseModel
	if err := json.Unmarshal(response.Data, &infoModel); err != nil {
		return errors.WithStack(err)
	}

	// download
	resp, err := http.Get(infoModel.Data.ExpiringDownloadURL)
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
