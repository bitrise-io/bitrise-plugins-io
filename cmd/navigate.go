package cmd

import (
	"github.com/bitrise-io/bitrise-plugins-io/views"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// navigateCmd represents the navigate command
var navigateCmd = &cobra.Command{
	Use:   "navigate",
	Short: "Interactive, navigation mode",
	Long: `Browe apps, builds, artifacts, logs, ...
in an interactive terminal/command line based UI.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		navigator := views.NewNavigationController()

		appsListView, err := views.NewAppsListViewController(navigator)
		if err != nil {
			return errors.WithStack(err)
		}

		navigator.PushViewController(appsListView)

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

		if err := app.SetRoot(navigator.View(), true).Run(); err != nil {
			panic(err)
		}

		return nil
	},
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
