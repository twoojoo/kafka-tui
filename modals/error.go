package modals

import (
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/types"
)

func ShowErrorModal(ui *types.UI, error string) {
	focused := ui.App.GetFocus()

	modal := tview.NewModal()
	ui.CentralView.AddItem(modal, 0, 0, false)

	modal.SetBackgroundColor(ui.Theme.Background)

	modal.SetText("Error: " + error)
	modal.SetTextColor(ui.Theme.ErrorColor)
	modal.AddButtons([]string{"Ok"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		ui.CentralView.RemoveItem(modal)
		ui.App.SetFocus(focused)
	})

	ui.App.SetFocus(modal)
}
