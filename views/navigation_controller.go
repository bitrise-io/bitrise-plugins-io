package views

import (
	"fmt"

	"github.com/rivo/tview"
)

// NavigationController ...
type NavigationController struct {
	pages *tview.Pages
	app   *tview.Application
}

// NewNavigationController ...
func NewNavigationController(app *tview.Application) *NavigationController {
	return &NavigationController{
		pages: tview.NewPages(),
		app:   app,
	}
}

// View ...
func (n *NavigationController) View() tview.Primitive {
	return n.pages
}

// PushViewController ...
func (n *NavigationController) PushViewController(vc ViewController) {
	pageID := fmt.Sprintf("page-%d", n.pages.GetPageCount())
	n.pages.AddAndSwitchToPage(pageID, vc.View(), true)
}

// PopViewController ...
func (n *NavigationController) PopViewController() {
	pageID := fmt.Sprintf("page-%d", n.pages.GetPageCount()-1)
	n.pages.RemovePage(pageID)
}

// FocusOnView ..
func (n *NavigationController) FocusOnView(view tview.Primitive) {
	n.app.SetFocus(view)
}
