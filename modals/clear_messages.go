package modals

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/types"
)

func ShowClearMessagesModal(ui *types.UI, topic string) {
	modal := tview.NewModal()
	ui.CentralView.AddItem(modal, 0, 0, false)

	modal.SetBackgroundColor(tcell.ColorDarkOrange)

	modal.SetText("Clear " + topic + " messages?\n(only if cleanup policy is set to 'delete')")
	modal.SetTextColor(tcell.ColorWhite)
	modal.AddButtons([]string{"NO", "YES"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 1 {
			clearMessages(ui, topic)
		}

		ui.App.SetFocus(ui.TopicsTable.Table)
		ui.CentralView.RemoveItem(modal)
	})

	ui.App.SetFocus(modal)
}

func clearMessages(ui *types.UI, topic string) {

}
