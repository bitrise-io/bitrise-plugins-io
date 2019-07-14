package cmd

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/colorstring"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(apps())
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.Flags().StringVar(&appsNextFlag, "next", "", "Next parameter for paging")
	appsCmd.Flags().StringVarP(&appsLimitFlag, "limit", "l", "", "Limit parameter for paging")
	appsCmd.Flags().StringVar(&appsSortFlag, "sort", string(services.SortAppsByLastBuildAt),
		fmt.Sprintf("Sort by parameter for listing. Options: [%s, %s]", services.SortAppsByLastBuildAt, services.SortAppsByCreatedAt))
}

type appsFormatter struct {
	*services.AppsListResponseModel
}

// Pretty ...
func (respModel *appsFormatter) Pretty() string {
	s := ""
	for _, aAppData := range respModel.Data {
		s += fmt.Sprintf("%s / %s (%s)\n", aAppData.Owner.Name, colorstring.Green(aAppData.Title), aAppData.Slug)
	}
	return s
}

func apps() error {
	response, err := services.GetBitriseAppsForUser(appsNextFlag, appsLimitFlag, services.AppSortBy(appsSortFlag))
	if err != nil {
		return errors.WithStack(err)
	}

	if response.Error != "" {
		return services.NewRequestFailedError(response)
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &appsFormatter{}))
}
