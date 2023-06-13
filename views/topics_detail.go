package views

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/kafka"
	"github.com/twoojoo/ktui/types"
	"github.com/twoojoo/ktui/utils"
)

func ShowTopicDetail(ui *types.UI, topic string) {
	topic = strings.Trim(topic, " ")
	info := ui.Topics[topic]
	metadata := kafka.GetTopicMetadata(ui, topic)

	// ui.TopicDetail.Clear()
	ui.TopicDetail = tview.NewFlex()
	ui.TopicDetail.SetBorder(true)
	ui.TopicDetail.SetBackgroundColor(ui.Theme.Background)
	ui.TopicDetail.SetDirection(0)
	ui.TopicDetail.SetTitle(" Topic detail: ")
	ui.TopicDetail.SetTitleAlign(0)

	// ui.CentralView.Clear()
	// ui.CentralView.AddItem(ui.TopicsTable.Container, 0, 1, true)
	ui.MainView.Clear()
	ui.MainView.AddItem(ui.SidePane, 20, 0, true)
	ui.MainView.AddItem(ui.CentralView, 0, 1, false)
	ui.MainView.AddItem(ui.TopicDetail, 60, 1, false)

	detailTitle := tview.NewTextView()
	detailTitle.SetText("\n " + topic)
	detailTitle.SetBackgroundColor(ui.Theme.Background)
	detailTitle.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
	detailTitle.SetTextColor(ui.Theme.PrimaryColor)

	kindText := utils.BuildDetailText(ui, " - internal: "+strconv.FormatBool(metadata.IsInternal))
	sizeText := utils.BuildDetailText(ui, " - size: "+utils.BytesToString(ui.TopicsSize[topic]))
	partitionsText := utils.BuildDetailText(ui, " - partitions: "+strconv.Itoa(len(metadata.Partitions)))
	messagesText := utils.BuildDetailText(ui, " - messages: ")
	replicaText := utils.BuildDetailText(ui, " - rep. factor: "+strconv.Itoa(int(info.ReplicationFactor)))

	filler := tview.NewTextView()
	filler.SetTextColor(ui.Theme.PrimaryColor)
	filler.SetBackgroundColor(ui.Theme.Background)

	ui.TopicDetail.AddItem(detailTitle, 3, 0, false)

	mainSubtitle := tview.NewTextView()
	mainSubtitle.SetBackgroundColor(ui.Theme.Background)
	mainSubtitle.SetText(" Info:")
	mainSubtitle.SetTextColor(ui.Theme.Foreground)
	ui.TopicDetail.AddItem(mainSubtitle, 1, 0, false)

	ui.TopicDetail.AddItem(kindText, 1, 0, false)
	ui.TopicDetail.AddItem(partitionsText, 1, 0, false)
	ui.TopicDetail.AddItem(replicaText, 1, 0, false)
	ui.TopicDetail.AddItem(messagesText, 1, 0, false)
	ui.TopicDetail.AddItem(sizeText, 1, 0, false)
	ui.TopicDetail.AddItem(utils.BuildDetailText(ui, ""), 1, 0, false)

	sortedConfigNames := utils.SortMapKeys(info.ConfigEntries)

	cfgSubtitle := tview.NewTextView()
	cfgSubtitle.SetBackgroundColor(ui.Theme.Background)
	cfgSubtitle.SetText(" Config:")
	cfgSubtitle.SetTextColor(ui.Theme.Foreground)
	ui.TopicDetail.AddItem(cfgSubtitle, 1, 0, false)

	for _, name := range sortedConfigNames {
		text := utils.BuildDetailText(ui, " - "+name+": "+*info.ConfigEntries[name])
		ui.TopicDetail.AddItem(text, 1, 0, false)
	}

	ui.TopicDetail.AddItem(filler, 0, 1, false)
}