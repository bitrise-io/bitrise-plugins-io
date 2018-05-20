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
	nextFlag  string
	limitFlag string
	sortFlag  string
)

const (
	sortFlagDefault = "last_build_at"
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
	appsCmd.Flags().StringVar(&nextFlag, "next", "", "Next parameter for paging")
	appsCmd.Flags().StringVarP(&limitFlag, "limit", "l", "", "Limit parameter for paging")
	appsCmd.Flags().StringVar(&sortFlag, "sort", sortFlagDefault, "Sort by parameter for listing. Options: [created_at, last_build_at]")
}

// AppsResponseModel ...
type AppsResponseModel struct {
	Data []struct {
		Title string `json:"title"`
		Slug  string `json:"slug"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"data"`
}

// Pretty ...
func (respModel *AppsResponseModel) Pretty() string {
	s := ""
	for _, aAppData := range respModel.Data {
		s += fmt.Sprintf("%s / %s (%s)\n", aAppData.Owner.Name, colorstring.Green(aAppData.Title), aAppData.Slug)
	}
	return s
}

func apps() error {
	// for some reason the Cobra Flag default doesn't seem to work for the sortFlag var
	// so doing it the manual way, ensuring default is set if needed
	if sortFlag == "" {
		sortFlag = sortFlagDefault
	}

	params := map[string]string{
		"next":    nextFlag,
		"limit":   limitFlag,
		"sort_by": sortFlag,
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
	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &AppsResponseModel{}))
}
