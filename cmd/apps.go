package cmd

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	appsNextFlag  string
	appsLimitFlag string
	appsSortFlag  string
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Get apps for user",
	Run: func(cmd *cobra.Command, args []string) {
		if err := apps(); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.Flags().StringVar(&appsNextFlag, "next", "", "Next parameter for paging")
	appsCmd.Flags().StringVarP(&appsLimitFlag, "limit", "l", "", "Limit parameter for paging")
	appsCmd.Flags().StringVar(&appsSortFlag, "sort", "last_build_at", "Sort by parameter for listing. Options: [created_at, last_build_at]")
}

// AppsListResponseModel ...
type AppsListResponseModel struct {
	Data []struct {
		Title string `json:"title"`
		Slug  string `json:"slug"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"data"`
}

// Pretty ...
func (respModel *AppsListResponseModel) Pretty() string {
	s := ""
	for _, aAppData := range respModel.Data {
		s += fmt.Sprintf("%s / %s (%s)\n", aAppData.Owner.Name, colorstring.Green(aAppData.Title), aAppData.Slug)
	}
	return s
}

func apps() error {
	params := map[string]string{
		"next":    appsNextFlag,
		"limit":   appsLimitFlag,
		"sort_by": appsSortFlag,
	}

	response, err := services.GetBitriseAppsForUser(params)
	if err != nil {
		return errors.WithStack(err)
	}

	if response.Error != "" {
		printErrorOutput(response.Error, formatFlag != "json")
		os.Exit(1)
		return nil
	}
	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &AppsListResponseModel{}))
}
