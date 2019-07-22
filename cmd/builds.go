package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/bitrise-io/bitrise-plugins-io/views/formatter"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	buildsAppSlugFlag string
	buildsSortFlag    string
	buildsNextFlag    string
	buildsLimitFlag   string
)

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Get builds for app",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(builds())
	},
}

func init() {
	rootCmd.AddCommand(buildsCmd)
	buildsCmd.Flags().StringVar(&buildsNextFlag, "next", "", "Next parameter for paging")
	buildsCmd.Flags().StringVarP(&buildsLimitFlag, "limit", "l", "", "Limit parameter for paging")
	buildsCmd.Flags().StringVar(&buildsSortFlag, "sort", "created_at", "Sort by parameter for listing. Options: [created_at, running_first]")
	buildsCmd.Flags().StringVarP(&buildsAppSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
}

type buildsFormatter struct {
	*services.BuildsListItemReponseModel
}
type buildsListFormatter struct {
	*services.BuildsListReponseModel
}

func coloredStatusText(statusText string) string {
	colorFN := colorstring.Blue
	switch statusText {
	case "success":
		colorFN = colorstring.Green
	case "error":
		colorFN = colorstring.Red
	case "aborted":
		colorFN = colorstring.Yellow
	}
	return colorFN(statusText)
}

func prettyPR(prID int) string {
	if prID == 0 {
		// empty string if zero
		return ""
	}
	return fmt.Sprintf("%d", prID)
}

// Pretty ...
func (respModel *buildsListFormatter) Pretty() string {
	buf := bytes.NewBuffer([]byte{})
	prettyTabWriter := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintln(prettyTabWriter, "#\t"+colorstring.Blue("Status")+"\tSlug\tTrigger Info\tMessage\tWorkflow"); err != nil {
		panic(err)
	}
	for _, aItem := range respModel.Data {
		fields := []string{
			fmt.Sprintf("%d", aItem.BuildNumber),
			coloredStatusText(aItem.StatusText),
			aItem.Slug,
			aItem.TriggerInfoString(),
			formatter.PrettyOneLinerText(aItem.CommitMessage),
			aItem.TriggeredWorkflow,
		}
		if _, err := fmt.Fprintln(prettyTabWriter, strings.Join(fields, "\t")); err != nil {
			panic(err)
		}
	}
	if err := prettyTabWriter.Flush(); err != nil {
		panic(err)
	}
	return buf.String()
}

func builds() error {
	params := map[string]string{
		"next":    buildsNextFlag,
		"limit":   buildsLimitFlag,
		"sort_by": buildsSortFlag,
	}

	response, err := services.GetBitriseBuildsForApp(buildsAppSlugFlag, params)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return services.NewRequestFailedError(response)
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &buildsListFormatter{}))
}
