package views

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/kafka"
	"github.com/twoojoo/ktui/modals"
	"github.com/twoojoo/ktui/types"
	"github.com/twoojoo/ktui/utils"
)

func ShowTopicsView(ui *types.UI) {
	topics := kafka.GetTopics(ui.AdminClient)
	ui.Topics = topics
	ui.TopicsMetadata = kafka.GetTopicsMetadata(ui.AdminClient, topics)

	topicNames := []string{}
	for topic := range ui.Topics {
		topicNames = append(topicNames, topic)
	}

	sort.Strings(topicNames)
	ui.TopicsSize = kafka.GetTopicsSize(ui.AdminClient, topicNames)

	ui.CentralView.SetBorder(false)

	ui.TopicsTable.Table.Clear()

	ui.TopicsTable.SetColumnNames([]string{
		" Name   ",
		" Internal   ",
		" Partitions   ",
		// " Out of sync replicas   ",
		" Replication factor   ",
		" Messages   ",
		" Size   ",
	}, ui.Theme.PrimaryColor)

	i := 1
	for _, topic := range topicNames {
		info := ui.Topics[topic]
		partitions := strconv.FormatInt(int64(info.NumPartitions), 10)
		repFactor := strconv.Itoa(int(info.ReplicationFactor))

		var isInternal string
		for _, item := range ui.TopicsMetadata {
			if item.Name == topic {
				isInternal = strconv.FormatBool(item.IsInternal)
				break
			}
		}

		ui.TopicsTable.Table.SetCell(i, 0, tview.NewTableCell(" "+topic+"   ").SetTextColor(ui.Theme.InEvidenceColor))
		ui.TopicsTable.Table.SetCell(i, 1, tview.NewTableCell(" "+isInternal+"   ").SetTextColor(ui.Theme.Foreground).SetStyle(tcell.StyleDefault.Attributes(tcell.AttrDim)))
		ui.TopicsTable.Table.SetCell(i, 2, tview.NewTableCell(" "+partitions+"   ").SetTextColor(ui.Theme.Foreground).SetStyle(tcell.StyleDefault.Attributes(tcell.AttrDim)))
		ui.TopicsTable.Table.SetCell(i, 3, tview.NewTableCell(" "+repFactor+"   ").SetTextColor(ui.Theme.Foreground).SetStyle(tcell.StyleDefault.Attributes(tcell.AttrDim)))
		ui.TopicsTable.Table.SetCell(i, 5, tview.NewTableCell(" "+utils.BytesToString(ui.TopicsSize[topic])+"   ").SetTextColor(ui.Theme.Foreground).SetStyle(tcell.StyleDefault.Attributes(tcell.AttrDim)))
		i++
	}

	ui.TopicsTable.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		rune := event.Rune()

		if rune == 'a' {
			modals.ShowAddTopicModal(ui)
		} else {
			row, _ := ui.TopicsTable.Table.GetSelection()
			topic := strings.Trim(ui.TopicsTable.Table.GetCell(row, 0).Text, " ")

			isInternal := false
			for _, item := range ui.TopicsMetadata {
				if item.Name == topic {
					isInternal = item.IsInternal
					break
				}
			}

			if isInternal && (rune == 'e' || rune == 'c' || rune == 'd') {
				modals.ShowErrorModal(ui, "can't edit internal topics")
				return event
			}

			if rune == 'e' {
				modals.ShowEditTopicModal(ui, topic)
			} else if rune == 'c' {
				modals.ShowClearMessagesModal(ui, topic)
			} else if rune == 'r' {

			} else if rune == 'd' {
				modals.ShowDeleteTopicModal(ui, topic)
			}
		}
		return event
	})
}
