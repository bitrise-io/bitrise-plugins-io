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
	view *tview.List
}

// View ...
func (vc *AppsListViewController) View() tview.Primitive {
	return vc.view
}

// NewAppsListViewController ...
func NewAppsListViewController(navigationController *NavigationController) (*AppsListViewController, error) {
	appsListView := tview.NewList()
	appsListView.SetTitle("Apps").SetBorder(true)
	// appsListView := tview.NewTable().SetBorders(true)
	// appsListView.SetTitle("Apps").SetBorder(true)

	// appsListView.SetCell(0, 0, &tview.TableCell{Text: "Owner", Align: tview.AlignCenter, Color: tcell.ColorYellow})
	// appsListView.SetCell(0, 1, &tview.TableCell{Text: "Title", Align: tview.AlignCenter, Color: tcell.ColorYellow})
	// appsListView.SetCell(0, 2, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow})

	apps, err := getApps()
	if err != nil {
		return nil, errors.WithStack(err)
	}

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
		// appsListView.SetCell(i+1, 0, &tview.TableCell{Text: aAppData.Owner.Name, Color: tcell.ColorDefault})
		// appsListView.SetCell(i+1, 1, &tview.TableCell{Text: aAppData.Title, Color: tcell.ColorDefault})
		// appsListView.SetCell(i+1, 2, &tview.TableCell{Text: aAppData.Slug, Color: tcell.ColorDefault})
	}

	appsListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'o' {
			currItemIdx := appsListView.GetCurrentItem()
			currHighlightedApp := apps.Data[currItemIdx]
			if err := utils.OpenPageInBrowser(currHighlightedApp.Slug, ""); err != nil {
				infoPopupVC := NewInfoPopupViewController(err.Error(), navigationController)
				navigationController.PushViewController(infoPopupVC)
			}
		}
		return event
	})

	return &AppsListViewController{
		view: appsListView,
	}, nil
}

func getApps() (services.AppsListResponseModel, error) {
	appListResp := services.AppsListResponseModel{}

	response, err := services.GetBitriseAppsForUser(map[string]string{})
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
