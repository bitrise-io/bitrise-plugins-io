package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	appSlugFlag     string
	buildsSortFlag  string
	buildsNextFlag  string
	buildsLimitFlag string
)

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Get builds for app",
	Run: func(cmd *cobra.Command, args []string) {
		if err := builds(); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildsCmd)
	buildsCmd.Flags().StringVar(&buildsNextFlag, "next", "", "Next parameter for paging")
	buildsCmd.Flags().StringVarP(&buildsLimitFlag, "limit", "l", "", "Limit parameter for paging")
	buildsCmd.Flags().StringVar(&buildsSortFlag, "sort", "created_at", "Sort by parameter for listing. Options: [created_at, running_first]")
	buildsCmd.Flags().StringVarP(&appSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
}

// BuildsReponseModel ...
type BuildsReponseModel struct {
	Slug                    string `json:"slug"`
	Status                  int    `json:"status"`
	StatusText              string `json:"status_text"`
	IsOnHold                bool   `json:"is_on_hold"`
	BuildNumber             int    `json:"build_number"`
	Branch                  string `json:"branch"`
	PullRequestTargetBranch string `json:"pull_request_target_branch"`
	PullRequestID           int    `json:"pull_request_id"`
	CommitHash              string `json:"commit_hash"`
	CommitMessage           string `json:"commit_message"`
	TriggeredWorkflow       string `json:"triggered_workflow"`
	TriggeredBy             string `json:"triggered_by"`
}

// BuildsListReponseModel ...
type BuildsListReponseModel struct {
	Data []BuildsReponseModel `json:"data"`
}

// BranchPRInfoString ...
func (respModel *BuildsReponseModel) BranchPRInfoString() string {
	if respModel.PullRequestID > 0 {
		return fmt.Sprintf("(#%d) %s > %s", respModel.PullRequestID, respModel.Branch, respModel.PullRequestTargetBranch)
	}
	return respModel.Branch
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

func prettyOneLinerText(msg string) string {
	s := fmt.Sprintf("%.50s", msg)         // print the first X chars
	s = strings.Replace(s, "\n", "↲", -1)  // replace newlines
	s = strings.Replace(s, "\r", "↲", -1)  // replace newlines
	s = strings.Replace(s, "\t", "  ", -1) // replace tabs
	return s
}

func prettyPR(prID int) string {
	if prID == 0 {
		// empty string if zero
		return ""
	}
	return fmt.Sprintf("%d", prID)
}

// Pretty ...
func (respModel *BuildsListReponseModel) Pretty() string {
	buf := bytes.NewBuffer([]byte{})
	prettyTabWriter := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintln(prettyTabWriter, "#\t"+colorstring.Blue("Status")+"\tSlug\tBranch / PR Info\tMessage\tWorkflow"); err != nil {
		panic(err)
	}
	for _, aItem := range respModel.Data {
		fields := []string{
			fmt.Sprintf("%d", aItem.BuildNumber),
			coloredStatusText(aItem.StatusText),
			aItem.Slug,
			aItem.BranchPRInfoString(),
			prettyOneLinerText(aItem.CommitMessage),
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

	response, err := services.GetBitriseBuildsForApp(appSlugFlag, params)
	if err != nil {
		return err
	}

	if response.Error != "" {
		printErrorOutput(response.Error, formatFlag != "json")
		os.Exit(1)
		return nil
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &BuildsListReponseModel{}))
}
