package modals

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/types"
)

const DEFAULT_TOPIC_NAME = "new-topic"
const DEFAULT_PARTITIONS = "1"
const DEFAULT_MAX_MSG_BYTES = "1048576"
const DEFAULT_RETENTION_BYTES = "-1"
const DEFAULT_RETENTION_MS = "604900000"
const DEFAULT_SEGMENT_BYTES = "134217728"
const DEFAULT_SEGMENT_MS = "12096000000"

type NewTopicFields struct {
	Name            string
	Compression     string
	RepFactor       int
	Partitions      int
	MaxMessageBytes int
	RetetntionBytes int
	RetentionMs     int
	SegmentBytes    int
	SegmentMs       int
}

func ShowAddTopicModal(ui *types.UI) {
	ui.IsModalOpen = true
	ui.Modal = tview.NewFlex()
	form := tview.NewForm()
	ui.Modal.AddItem(form, 0, 1, true)
	ui.CentralView.AddItem(ui.Modal, 23, 0, false)

	ui.Modal.SetBorder(true)
	ui.Modal.SetTitle(" Add new topic ")
	ui.Modal.SetTitleColor(ui.Theme.InfoColor)
	ui.Modal.SetBorderColor(ui.Theme.InfoColor)
	ui.Modal.SetTitleAlign(0)
	ui.Modal.SetBackgroundColor(ui.Theme.Background)
	form.SetBackgroundColor(ui.Theme.Background)

	fields := NewTopicFields{}

	form.AddInputField(" Name: ", DEFAULT_TOPIC_NAME, 20, func(text string, _ rune) bool {
		fields.Name = text
		return true
	}, nil)

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
		exitModal(ui)
	})

	form.AddButton("Add", func() {
		createTopic(ui, fields)
		exitModal(ui)
	})

	form.SetButtonBackgroundColor(ui.Theme.Background)
	form.SetButtonTextColor(ui.Theme.Foreground)
	form.SetButtonActivatedStyle(tcell.StyleDefault.Background(ui.Theme.Background).Foreground(ui.Theme.PrimaryColor))

	ui.Modal.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEscape {
			exitModal(ui)
		}

		return e
	})

	ui.App.SetFocus(ui.Modal)
}

func exitModal(ui *types.UI) {
	ui.IsModalOpen = false
	ui.CentralView.Clear()
	ui.CentralView.AddItem(ui.TopicsTable.Container, 0, 1, false)
	ui.App.SetFocus(ui.TopicsTable.Table)
}

func createTopic(ui *types.UI, fields NewTopicFields) {

}
