package modals

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/types"
)

type EditTopicFields struct {
	RepFactor       int
	Partitions      int
	MaxMessageBytes int
	RetetntionBytes int
	RetentionMs     int
	SegmentBytes    int
	SegmentMs       int
}

func ShowEditTopicModal(ui *types.UI, topic string) {
	ui.IsModalOpen = true
	ui.Modal = tview.NewFlex()
	form := tview.NewForm()
	ui.Modal.AddItem(form, 0, 1, true)
	ui.CentralView.AddItem(ui.Modal, 21, 1, false)

	ui.Modal.SetBorder(true)
	ui.Modal.SetTitle(" Edit " + topic + " ")
	ui.Modal.SetTitleColor(tcell.ColorYellow)
	ui.Modal.SetBorderColor(tcell.ColorYellow)
	ui.Modal.SetTitleAlign(0)
	ui.Modal.SetBackgroundColor(ui.Theme.Background)
	form.SetBackgroundColor(ui.Theme.Background)

	fields := NewTopicFields{}

	form.AddInputField(" Partitions: ", DEFAULT_PARTITIONS, 20, func(text string, lastChar rune) bool {
		num, err := strconv.Atoi(text)

		if err != nil {

			return false
		}

		fields.Partitions = num
		return true
	}, nil)

	form.AddInputField(" Rep. Factor: ", "1", 20, func(text string, lastChar rune) bool {
		num, err := strconv.Atoi(text)

		if err != nil {
			return false
		}

		fields.RepFactor = num
		return true
	}, nil)

	form.AddInputField(" compression: ", "", 20, func(text string, lastChar rune) bool {
		fields.Compression = text
		return true
	}, nil)

	form.AddInputField(" max.message.bytes: ", DEFAULT_MAX_MSG_BYTES, 20, func(text string, lastChar rune) bool {
		num, err := strconv.Atoi(text)

		if err != nil {
			return false
		}

		fields.RepFactor = num
		return true
	}, nil)

	form.AddInputField(" retention.bytes: ", DEFAULT_RETENTION_BYTES, 20, func(text string, lastChar rune) bool {
		num, err := strconv.Atoi(text)

		if err != nil {
			return false
		}

		fields.RepFactor = num
		return true
	}, nil)

	form.AddInputField(" retention.ms: ", DEFAULT_RETENTION_MS, 20, func(text string, lastChar rune) bool {
		num, err := strconv.Atoi(text)

		if err != nil {
			return false
		}

		fields.RepFactor = num
		return true
	}, nil)

	form.AddInputField(" segment.bytes: ", DEFAULT_SEGMENT_BYTES, 20, func(text string, lastChar rune) bool {
		num, err := strconv.Atoi(text)

		if err != nil {
			return false
		}

		fields.RepFactor = num
		return true
	}, nil)

	form.AddInputField(" segment.ms: ", DEFAULT_SEGMENT_MS, 20, func(text string, lastChar rune) bool {
		num, err := strconv.Atoi(text)

		if err != nil {
			return false
		}

		fields.RepFactor = num
		return true
	}, nil)

	form.SetFieldBackgroundColor(ui.Theme.Background)
	form.SetFieldTextColor(ui.Theme.InEvidenceColor)
	form.SetLabelColor(ui.Theme.Foreground)

	form.AddButton("Cancel", func() {
		exitEditTopicModal(ui)
	})

	form.AddButton("Edit", func() {
		createTopic(ui, fields)
		exitEditTopicModal(ui)
	})

	form.SetButtonBackgroundColor(ui.Theme.Background)
	form.SetButtonTextColor(ui.Theme.Foreground)
	form.SetButtonActivatedStyle(tcell.StyleDefault.Background(ui.Theme.Background).Foreground(ui.Theme.PrimaryColor))

	ui.Modal.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEscape {
			exitEditTopicModal(ui)
		}

		return e
	})

	ui.App.SetFocus(ui.Modal)
}

func exitEditTopicModal(ui *types.UI) {
	ui.IsModalOpen = false
	ui.CentralView.Clear()
	ui.CentralView.AddItem(ui.TopicsTable.Container, 0, 1, false)
	ui.App.SetFocus(ui.TopicsTable.Table)
}

func editTopic(ui *types.UI, fields NewTopicFields) {

}
