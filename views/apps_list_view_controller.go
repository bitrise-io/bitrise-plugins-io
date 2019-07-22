package views

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/bitrise-io/bitrise-plugins-io/utils"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

// AppsListViewController ...
type AppsListViewController struct {
	view                 tview.Primitive
	apps                 services.AppsListResponseModel
	titleFilter          string
	appsListView         *tview.List
	navigationController *NavigationController
}

// View ...
func (vc *AppsListViewController) View() tview.Primitive {
	return vc.view
}

func xreloadAppsListView(appsListView *tview.List, apps services.AppsListResponseModel, navigationController *NavigationController) {
	appsListView.Clear()
	for i, aAppData := range apps.Data {
		aAppSlug := aAppData.Slug
		appsListView.AddItem(
			fmt.Sprintf("[#%d] %s / %s (%s)", i, aAppData.Owner.Name, aAppData.Title, aAppData.Slug),
			"",
			0,
			func() {
				buildsListViewController, err := NewBuildsListViewController(aAppSlug, navigationController)
				if err != nil {
					panic(err)
				}
				navigationController.PushViewController(buildsListViewController)
			})
	}
}

func xreloadAppsDataAndView(titleFilter string, appsListView *tview.List, navigationController *NavigationController) (services.AppsListResponseModel, error) {
	apps, err := getApps(titleFilter)
	if err != nil {
		return apps, errors.WithStack(err)
	}
	xreloadAppsListView(appsListView, apps, navigationController)

	timeNowStr := time.Now().Format("15:04:05")
	listTitle := fmt.Sprintf(" Apps [%s] ", timeNowStr)
	if len(titleFilter) > 0 {
		listTitle = fmt.Sprintf(" Apps (%s) [%s] ", titleFilter, timeNowStr)
	}
	appsListView.SetTitle(listTitle)
	return apps, nil
}

func (vc *AppsListViewController) reloadAppsDataAndView() error {
	apps, err := xreloadAppsDataAndView(vc.titleFilter, vc.appsListView, vc.navigationController)
	if err != nil {
		return errors.WithStack(err)
	}
	vc.apps = apps
	return nil
}

func (vc *AppsListViewController) getSelectedAppData() services.AppsListItemResponseModel {
	currItemIdx := vc.appsListView.GetCurrentItem()
	return vc.apps.Data[currItemIdx]
}

// NewAppsListViewController ...
func NewAppsListViewController(navigationController *NavigationController) (*AppsListViewController, error) {
	layoutView := tview.NewFlex().SetDirection(tview.FlexRow)

	appsListView := tview.NewList()
	appsListView.SetTitle("Apps").SetBorder(true)
	appsListView.SetSelectedBackgroundColor(selectedBackgroundColor)

	appsListViewController := &AppsListViewController{
		navigationController: navigationController,
		view:                 layoutView,
		appsListView:         appsListView,
	}

	if err := appsListViewController.reloadAppsDataAndView(); err != nil {
		return nil, errors.WithStack(err)
	}

	searchInputField := tview.NewInputField().
		SetLabel("Search: ")
	searchInputField.SetLabelColor(tcell.ColorWhiteSmoke)
	searchInputField.SetFieldBackgroundColor(tcell.ColorMediumPurple)
	searchInputField.SetDoneFunc(func(key tcell.Key) {
		layoutView.RemoveItem(searchInputField)
		navigationController.FocusOnView(layoutView)
		// reload apps list with the new filter
		appsListViewController.titleFilter = searchInputField.GetText()
		if err := appsListViewController.reloadAppsDataAndView(); err != nil {
			panic(errors.WithStack(err))
		}
	})

	appsListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEscape:
			navigationController.FocusOnView(appsListView)
		case event.Rune() == 'o':
			currHighlightedApp := appsListViewController.getSelectedAppData()
			if err := utils.OpenPageInBrowser(currHighlightedApp.Slug, ""); err != nil {
				infoPopupVC := NewInfoPopupViewController(err.Error(), navigationController)
				navigationController.PushViewController(infoPopupVC)
			}
		case event.Rune() == 'f':
			layoutView.AddItem(searchInputField, 0, 1, true)
			navigationController.FocusOnView(searchInputField)
		case event.Rune() == 'r':
			if err := appsListViewController.reloadAppsDataAndView(); err != nil {
				panic(errors.WithStack(err))
			}
		case event.Rune() == '?':
			hotkeysPopupVC := NewInfoPopupViewController(`Hotkeys:

Base navigation:
- arrow keys: move item selection
- enter: select highlighted item
- esc: back

Other:
- f: filter
- o: open in browser
- r: reload
- ctrl+q: quit
- ctrl+c: quit
- ?: show hotkeys
`, navigationController)
			navigationController.PushViewController(hotkeysPopupVC)
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
				"enter":      "select app",
				"esc":        "back",
				"?":          "show help/hotkeys",
				"ctrl+q":     "quit",
			},
			map[string]string{
				"f": "filter",
				"o": "open in browser",
				"r": "reload apps",
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
	layoutView.AddItem(appsListView, 0, 1, true)

	return appsListViewController, nil
}

func getApps(titleFilter string) (services.AppsListResponseModel, error) {
	appListResp := services.AppsListResponseModel{}

	response, err := services.GetBitriseAppsForUser("", "", services.SortAppsByLastBuildAt, titleFilter)
	if err != nil {
		return appListResp, errors.WithStack(err)
	}

	if response.Error != "" {
		return appListResp, services.NewRequestFailedError(response)
	}

	if err := json.Unmarshal(response.Data, &appListResp); err != nil {
		return appListResp, errors.WithStack(err)
	}

	return appListResp, nil
}
