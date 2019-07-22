package views

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

// LogsViewController ...
type LogsViewController struct {
	appSlug   string
	buildSlug string
	view      tview.Primitive
}

// View ...
func (vc *LogsViewController) View() tview.Primitive {
	return vc.view
}

// NewLogsViewController ...
func NewLogsViewController(appSlug, buildSlug string, navigationController *NavigationController) (*LogsViewController, error) {
	logView := tview.NewTextView()
	logView.SetBorder(true)
	logView.SetDynamicColors(true)
	logView.SetTitle(fmt.Sprintf("Logs of (app: %s) (build: %s)", appSlug, buildSlug))
	logView.SetDoneFunc(func(key tcell.Key) {
		navigationController.PopViewController()
	})

	layoutView := tview.NewFlex().SetDirection(tview.FlexRow)
	logsViewController := &LogsViewController{
		appSlug:   appSlug,
		buildSlug: buildSlug,
		view:      layoutView,
	}

	layoutView.AddItem(logView, 0, 1, true)

	buildLogText, err := loadBuildLog(appSlug, buildSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ansiLogViewWriter := tview.ANSIWriter(logView)
	fmt.Fprintf(ansiLogViewWriter, "%s", buildLogText)

	return logsViewController, nil
}

func loadBuildLog(appSlug, buildSlug string) (string, error) {
	params := map[string]string{}
	response, err := services.GetBuildLogInfo(appSlug, buildSlug, params)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if response.Error != "" {
		return "", errors.WithStack(services.NewRequestFailedError(response))
	}

	logInfo := struct {
		ExpiringRawLogURL string `json:"expiring_raw_log_url"`
		IsArchived        bool   `json:"is_archived"`
	}{}

	if err := json.Unmarshal(response.Data, &logInfo); err != nil {
		return "", errors.WithStack(err)
	}

	fullLogData, err := services.LoadFullLog(logInfo.ExpiringRawLogURL)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(fullLogData), nil
}
