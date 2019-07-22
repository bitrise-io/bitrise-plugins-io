package views

import "github.com/rivo/tview"

// InfoPopupViewController ...
type InfoPopupViewController struct {
	view *tview.Modal
}

// View ...
func (vc *InfoPopupViewController) View() tview.Primitive {
	return vc.view
}

// NewInfoPopupViewController ...
func NewInfoPopupViewController(infoText string, navigationController *NavigationController) *InfoPopupViewController {
	modal := tview.NewModal().
		SetText(infoText).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				navigationController.PopViewController()
			}
		})

	return &InfoPopupViewController{
		view: modal,
	}
}
