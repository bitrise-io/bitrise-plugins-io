package views

import "github.com/rivo/tview"

// ViewController ...
type ViewController interface {
	View() tview.Primitive
}
