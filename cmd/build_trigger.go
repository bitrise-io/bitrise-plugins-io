package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	buildTriggerParamsFlag string
)

// buildTriggerCmd represents the trigger command
var buildTriggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger build",
	Long:  `Trigger a new build`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.WithStack(triggerBuild())
	},
}

func init() {
	buildsCmd.AddCommand(buildTriggerCmd)
	buildTriggerCmd.Flags().StringVarP(&buildsAppSlugFlag, "app", "a", "", "Slug of the app where the builds belong to")
	buildTriggerCmd.Flags().StringVar(&buildTriggerParamsFlag, "params", "", "Trigger parameters (in JSON format)")
}

// BuildTriggerResponseModel ...
type BuildTriggerResponseModel struct {
	Status            string `json:"status"`
	Message           string `json:"message"`
	Slug              string `json:"slug"`
	Service           string `json:"service"`
	BuildSlug         string `json:"build_slug"`
	BuildNumber       int    `json:"build_number"`
	BuildURL          string `json:"build_url"`
	TriggeredWorkflow string `json:"triggered_workflow"`
}

// Pretty ...
func (respModel *BuildTriggerResponseModel) Pretty() string {
	buf := bytes.NewBuffer([]byte{})
	prettyTabWriter := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintln(prettyTabWriter, "#\t"+colorstring.Blue("Status")+"\tMessage\tBuild Slug\tBuild URL\tWorkflow"); err != nil {
		panic(err)
	}
	var statusText string
	if respModel.Status == "ok" {
		statusText = colorstring.Green(respModel.Status)
	} else {
		statusText = colorstring.Red(respModel.Status)
	}
	fields := []string{
		fmt.Sprintf("%d", respModel.BuildNumber),
		statusText,
		respModel.Message,
		respModel.BuildSlug,
		respModel.BuildURL,
		respModel.TriggeredWorkflow,
	}
	if _, err := fmt.Fprintln(prettyTabWriter, strings.Join(fields, "\t")); err != nil {
		panic(err)
	}

	if err := prettyTabWriter.Flush(); err != nil {
		panic(err)
	}
	return buf.String()
}

func triggerBuild() error {
	params := map[string]interface{}{}
	err := json.Unmarshal([]byte(buildTriggerParamsFlag), &params)
	if err != nil {
		return err
	}

	response, err := services.TriggerBitriseBuildForApp(buildsAppSlugFlag, params)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return NewRequestFailedError(response)
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != formatJSON, &BuildTriggerResponseModel{}))
}
