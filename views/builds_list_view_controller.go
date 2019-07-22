package views

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/bitrise-io/bitrise-plugins-io/utils"
	"github.com/bitrise-io/bitrise-plugins-io/views/formatter"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

// BuildsListViewController ...
type BuildsListViewController struct {
	appSlug              string
	view                 tview.Primitive
	buildsListView       *tview.List
	navigationController *NavigationController
	builds               services.BuildsListReponseModel
}

// View ...
func (vc *BuildsListViewController) View() tview.Primitive {
	return vc.view
}

func (vc *BuildsListViewController) reloadBuilds() error {
	buildListResp, err := getBuilds(vc.appSlug)
	if err != nil {
		return errors.WithStack(err)
	}
	vc.builds = buildListResp

	vc.buildsListView.Clear()
	for _, aBuildData := range buildListResp.Data {
		vc.buildsListView.AddItem(
			fmt.Sprintf("[#%d] %s (workflow:%s) (branch:%s) (at:%s)",
				aBuildData.BuildNumber,
				formatter.PrettyOneLinerText(aBuildData.CommitMessage),
				aBuildData.TriggeredWorkflow,
				aBuildData.Branch,
				aBuildData.TriggeredAt,
			),
			// fmt.Sprintf("[#%d] %+v", i, aBuildData),
			"",
			0,
			nil)
	}

	timeNowStr := time.Now().Format("15:04:05")
	vc.buildsListView.SetTitle(fmt.Sprintf(" Builds (%s) [%s] ", vc.appSlug, timeNowStr))

	return nil
}

func (vc *BuildsListViewController) getSelectedBuildData() services.BuildsListItemReponseModel {
	currItemIdx := vc.buildsListView.GetCurrentItem()
	return vc.builds.Data[currItemIdx]
}

func getBuilds(appSlug string) (services.BuildsListReponseModel, error) {
	buildListResp := services.BuildsListReponseModel{}

	response, err := services.GetBitriseBuildsForApp(appSlug, map[string]string{})
	if err != nil {
		return buildListResp, errors.WithStack(err)
	}

	if response.Error != "" {
		return buildListResp, errors.WithStack(services.NewRequestFailedError(response))
	}

	if err := json.Unmarshal(response.Data, &buildListResp); err != nil {
		return buildListResp, errors.WithStack(err)
	}

	return buildListResp, nil
}

// NewBuildsListViewController ...
func NewBuildsListViewController(appSlug string, navigationController *NavigationController) (*BuildsListViewController, error) {
	// item selected
	buildsListView := tview.NewList()
	buildsListView.SetBorder(true)
	buildsListView.SetSelectedBackgroundColor(selectedBackgroundColor)

	layoutView := tview.NewFlex().SetDirection(tview.FlexRow)
	buildsListViewController := &BuildsListViewController{
		appSlug:              appSlug,
		view:                 layoutView,
		buildsListView:       buildsListView,
		navigationController: navigationController,
	}

	if err := buildsListViewController.reloadBuilds(); err != nil {
		return nil, errors.WithStack(err)
	}

	buildsListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEscape:
			navigationController.PopViewController()
		case event.Rune() == 'o':
			currHighlightedBuild := buildsListViewController.getSelectedBuildData()
			if err := utils.OpenPageInBrowser(appSlug, currHighlightedBuild.Slug); err != nil {
				infoPopupVC := NewInfoPopupViewController(err.Error(), navigationController)
				navigationController.PushViewController(infoPopupVC)
			}
		case event.Rune() == 'r':
			if err := buildsListViewController.reloadBuilds(); err != nil {
				panic(errors.WithStack(err))
			}
		case event.Rune() == 'l':
			currHighlightedBuild := buildsListViewController.getSelectedBuildData()
			logsViewController, err := NewLogsViewController(appSlug, currHighlightedBuild.Slug, navigationController)
			if err != nil {
				panic(err)
			}
			navigationController.PushViewController(logsViewController)
		}
		return event
	})

	header := tview.NewFlex()
	{
		header.SetDirection(tview.FlexColumn)
		infosTable := tview.NewTable()
		header.AddItem(infosTable, 0, 1, false)

		commands := []map[string]string{
			map[string]string{
				"arrow keys": "move item selection",
				// "enter":      "select build",
				"esc": "back",
				// "?":          "show help/hotkeys",
				"ctrl+q": "quit",
			},
			map[string]string{
				// "f": "filter",
				"o": "open in browser",
				"r": "reload builds",
				"l": "logs",
			},
		}
		for cmdRowIdx, commandRow := range commands {
			row := 0
			for key, description := range commandRow {
				keyCell := tview.NewTableCell(key)
				keyCell.SetTextColor(tcell.ColorMediumPurple)
				infosTable.SetCell(row, (cmdRowIdx*2)+1, keyCell)

				descriptionCell := tview.NewTableCell(description)
				descriptionCell.SetTextColor(tcell.ColorWhiteSmoke)
				infosTable.SetCell(row, (cmdRowIdx*2)+0, descriptionCell)

				row++
			}
		}
	}

	layoutView.AddItem(header, 7, 1, false)
	layoutView.AddItem(buildsListView, 0, 1, true)

	return buildsListViewController, nil
}
