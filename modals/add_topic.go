package modals

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/types"
)

func ShowAddTopicModal(ui *types.UI) {
	ui.IsModalOpen = true
	ui.Modal = tview.NewForm()
	ui.CentralView.AddItem(ui.Modal, 0, 1, true)

	ui.Modal.SetBorder(true)
	ui.Modal.SetBackgroundColor(ui.Theme.Background)

	ui.Modal.AddInputField("Name: ", "new-topic", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})
	ui.Modal.AddInputField("Partitions: ", "1", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})
	ui.Modal.AddInputField("Rep. Factor: ", "1", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})
	ui.Modal.AddInputField("max.message.bytes: ", "-1", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})
	ui.Modal.AddInputField("retention.bytes: ", "-1", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})
	ui.Modal.AddInputField("retention.ms: ", "-1", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})
	ui.Modal.AddInputField("segment.bytes: ", "-1", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})
	ui.Modal.AddInputField("segment.ms: ", "-1", 20, func(text string, lastChar rune) bool { return true }, func(text string) {})

	ui.Modal.SetFieldBackgroundColor(ui.Theme.Background)
	ui.Modal.SetFieldTextColor(ui.Theme.PrimaryColor)

	ui.Modal.AddButton("Cancel", func() {
		exitModal(ui)
	})

	ui.Modal.AddButton("Ok", func() {
		exitModal(ui)
	})

	ui.Modal.SetButtonBackgroundColor(ui.Theme.Background)
	ui.Modal.SetButtonTextColor(ui.Theme.Foreground)
	ui.Modal.SetButtonActivatedStyle(tcell.StyleDefault.Background(ui.Theme.Background).Foreground(ui.Theme.PrimaryColor))

	ui.App.SetFocus(ui.Modal)
}

func exitModal(ui *types.UI) {
	ui.IsModalOpen = false
	ui.CentralView.Clear()
	ui.CentralView.AddItem(ui.TopicsTable.Container, 0, 1, false)
	ui.App.SetFocus(ui.TopicsTable.Table)
}
