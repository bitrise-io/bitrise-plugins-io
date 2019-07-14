package views

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/bitrise-io/bitrise-plugins-io/utils"
	"github.com/bitrise-io/bitrise-plugins-io/views/formatter"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

// BuildsListViewController ...
type BuildsListViewController struct {
	view *tview.List
}

// View ...
func (vc *BuildsListViewController) View() tview.Primitive {
	return vc.view
}

// NewBuildsListViewController ...
func NewBuildsListViewController(appSlug string, navigationController *NavigationController) (*BuildsListViewController, error) {
	// item selected
	buildsListView := tview.NewList()
	buildsListView.SetTitle(fmt.Sprintf("Builds (%s)", appSlug)).SetBorder(true)

	response, err := services.GetBitriseBuildsForApp(appSlug, map[string]string{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if response.Error != "" {
		return nil, errors.WithStack(services.NewRequestFailedError(response))
	}

	buildListResp := services.BuildsListReponseModel{}
	if err := json.Unmarshal(response.Data, &buildListResp); err != nil {
		panic(errors.WithStack(err))
	}

	for i, aBuildData := range buildListResp.Data {
		buildsListView.AddItem(
			fmt.Sprintf("[#%d] %d / %s (%s)", i, aBuildData.BuildNumber, formatter.PrettyOneLinerText(aBuildData.CommitMessage), aBuildData.TriggeredWorkflow),
			// fmt.Sprintf("[#%d] %+v", i, aBuildData),
			"",
			0,
			nil)
	}

	buildsListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			navigationController.PopViewController()
		}
		if event.Rune() == 'o' {
			currItemIdx := buildsListView.GetCurrentItem()
			currHighlightedBuild := buildListResp.Data[currItemIdx]
			if err := utils.OpenPageInBrowser(appSlug, currHighlightedBuild.Slug); err != nil {
				infoPopupVC := NewInfoPopupViewController(err.Error(), navigationController)
				navigationController.PushViewController(infoPopupVC)
			}
		}
		return event
	})

	return &BuildsListViewController{
		view: buildsListView,
	}, nil
}
