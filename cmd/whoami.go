package cmd

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Print info about the authenticated user",
	Long:  `Print info about the authenticated user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(whoami())
	},
}

func init() {
	authCmd.AddCommand(whoamiCmd)
}

// MeResponseModel ...
type MeResponseModel struct {
	Data struct {
		Username string `json:"username"`
		Slug     string `json:"slug"`
	} `json:"data"`
}

// Pretty ...
func (respModel *MeResponseModel) Pretty() string {
	return fmt.Sprintf("%s (%s)", respModel.Data.Username, respModel.Data.Slug)
}

func whoami() error {
	response, err := services.ValidateAuthToken()
	if err != nil {
		return err
	}

	if response.Error != "" {
		return services.NewRequestFailedError(response)
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &MeResponseModel{}))
}
