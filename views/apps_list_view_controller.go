package views

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/bitrise-io/bitrise-plugins-io/utils"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

// AppsListViewController ...
type AppsListViewController struct {
	view tview.Primitive
	apps services.AppsListResponseModel
}

// View ...
func (vc *AppsListViewController) View() tview.Primitive {
	return vc.view
}

func reloadAppsListView(appsListView *tview.List, apps services.AppsListResponseModel, navigationController *NavigationController) {
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

func reloadAppsDataAndView(titleFilter string, appsListView *tview.List, navigationController *NavigationController) (services.AppsListResponseModel, error) {
	apps, err := getApps(titleFilter)
	if err != nil {
		return apps, errors.WithStack(err)
	}
	reloadAppsListView(appsListView, apps, navigationController)

	listTitle := "Apps"
	if len(titleFilter) > 0 {
		listTitle = fmt.Sprintf("Apps (%s)", titleFilter)
	}
	appsListView.SetTitle(listTitle)
	return apps, nil
}

// NewAppsListViewController ...
func NewAppsListViewController(navigationController *NavigationController) (*AppsListViewController, error) {
	layoutView := tview.NewFlex().SetDirection(tview.FlexRow)
	appsListViewController := &AppsListViewController{
		view: layoutView,
	}

	appsListView := tview.NewList()
	appsListView.SetTitle("Apps").SetBorder(true)

	apps, err := reloadAppsDataAndView("", appsListView, navigationController)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	searchInputField := tview.NewInputField().
		SetLabel("Search: ")
	searchInputField.SetDoneFunc(func(key tcell.Key) {
		layoutView.RemoveItem(searchInputField)
		navigationController.FocusOnView(layoutView)
		// reload apps list with the new filter
		titleFilter := searchInputField.GetText()
		a, err := reloadAppsDataAndView(titleFilter, appsListView, navigationController)
		if err != nil {
			panic(errors.WithStack(err))
		}
		apps = a
	})

	appsListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEscape:
			navigationController.FocusOnView(appsListView)
		case event.Rune() == 'o':
			currItemIdx := appsListView.GetCurrentItem()
			currHighlightedApp := apps.Data[currItemIdx]
			if err := utils.OpenPageInBrowser(currHighlightedApp.Slug, ""); err != nil {
				infoPopupVC := NewInfoPopupViewController(err.Error(), navigationController)
				navigationController.PushViewController(infoPopupVC)
			}
		case event.Rune() == 'f':
			layoutView.AddItem(searchInputField, 0, 1, true)
			navigationController.FocusOnView(searchInputField)
		case event.Rune() == '?':
			hotkeysPopupVC := NewInfoPopupViewController(`Hotkeys:

Base navigation:
- arrow keys: move item selection
- enter: select highlighted item
- esc: back

Other:
- f: filter
- o: open in browser
- ctrl+q: quit
- ctrl+c: quit
`, navigationController)
			navigationController.PushViewController(hotkeysPopupVC)
		}
		return event
	})

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
