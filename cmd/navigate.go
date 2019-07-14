package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-io/services"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// navigateCmd represents the navigate command
var navigateCmd = &cobra.Command{
	Use:   "navigate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pages := tview.NewPages()

		appsListView := tview.NewList()
		appsListView.SetTitle("Apps").SetBorder(true)
		// appsListView := tview.NewTable().SetBorders(true)
		// appsListView.SetTitle("Apps").SetBorder(true)

		response, err := services.GetBitriseAppsForUser(map[string]string{})
		if err != nil {
			return errors.WithStack(err)
		}

		if response.Error != "" {
			return NewRequestFailedError(response)
		}

		appListResp := AppsListResponseModel{}
		if err := json.Unmarshal(response.Data, &appListResp); err != nil {
			return errors.WithStack(err)
		}
		// appsListView.SetCell(0, 0, &tview.TableCell{Text: "Owner", Align: tview.AlignCenter, Color: tcell.ColorYellow})
		// appsListView.SetCell(0, 1, &tview.TableCell{Text: "Title", Align: tview.AlignCenter, Color: tcell.ColorYellow})
		// appsListView.SetCell(0, 2, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow})
		for i, aAppData := range appListResp.Data {
			aAppSlug := aAppData.Slug
			appsListView.AddItem(
				fmt.Sprintf("[#%d] %s / %s (%s)", i, aAppData.Owner.Name, aAppData.Title, aAppData.Slug),
				"",
				0,
				func() {
					buildsListView := appSelected(aAppSlug, pages)
					pages.AddAndSwitchToPage("Builds", buildsListView, true)
				})
			// appsListView.SetCell(i+1, 0, &tview.TableCell{Text: aAppData.Owner.Name, Color: tcell.ColorDefault})
			// appsListView.SetCell(i+1, 1, &tview.TableCell{Text: aAppData.Title, Color: tcell.ColorDefault})
			// appsListView.SetCell(i+1, 2, &tview.TableCell{Text: aAppData.Slug, Color: tcell.ColorDefault})
		}

		pages.AddPage("Apps", appsListView, true, true)

		app := tview.NewApplication()
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlQ {
				app.Stop()
			}
			// if event.Rune() == 'q' {
			// 	app.Stop()
			// }
			return event
		})

		if err := app.SetRoot(pages, true).Run(); err != nil {
			panic(err)
		}

		return nil
	},
}

func appSelected(selectedAppSlug string, pages *tview.Pages) *tview.List {
	// item selected
	buildsListView := tview.NewList()
	buildsListView.SetTitle(fmt.Sprintf("Builds (%s)", selectedAppSlug)).SetBorder(true)

	response, err := services.GetBitriseBuildsForApp(selectedAppSlug, map[string]string{})
	if err != nil {
		panic(err)
	}

	if response.Error != "" {
		panic(NewRequestFailedError(response))
	}

	buildListResp := BuildsListReponseModel{}
	if err := json.Unmarshal(response.Data, &buildListResp); err != nil {
		panic(errors.WithStack(err))
	}

	for i, aBuildData := range buildListResp.Data {
		buildsListView.AddItem(
			fmt.Sprintf("[#%d] %d / %s (%s)", i, aBuildData.BuildNumber, prettyOneLinerText(aBuildData.CommitMessage), aBuildData.TriggeredWorkflow),
			// fmt.Sprintf("[#%d] %+v", i, aBuildData),
			"",
			0,
			nil)
	}

	buildsListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// pages.SwitchToPage("Apps")
			pages.RemovePage("Builds")
		}
		return event
	})

	return buildsListView
}

func init() {
	rootCmd.AddCommand(navigateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// navigateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// navigateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
