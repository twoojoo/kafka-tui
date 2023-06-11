package main

import (
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	adminClient                *sarama.ClusterAdmin
	theme                      *Theme
	app                        *tview.Application
	sidePane                   *tview.List
	view                       *tview.Grid
	consumersTable             *SearchableTable
	topicsTable                *SearchableTable
	topics                     map[string]sarama.TopicDetail
	topicsMetadata             []*sarama.TopicMetadata
	topicsSize                 map[string]int
	consumerGroups             map[string]string
	consumerGroupsOffsets      map[string]*sarama.OffsetFetchResponse
	consumerGroupsDescriptions []*sarama.GroupDescription
}

type Theme struct {
	Background      tcell.Color
	Foreground      tcell.Color
	PrimaryColor    tcell.Color
	InEvidenceColor tcell.Color
	InfoColor       tcell.Color
	ErrorColor      tcell.Color
	WarningColor    tcell.Color
}

func main() {
	app := tview.NewApplication()

	theme := Theme{
		Background:      tcell.ColorReset,
		Foreground:      tcell.ColorWhite,
		PrimaryColor:    tcell.ColorBlue,
		InEvidenceColor: tcell.ColorYellow,
		InfoColor:       tcell.ColorGreen,
		ErrorColor:      tcell.ColorRed,
		WarningColor:    tcell.ColorYellow,
	}

	sidePane := tview.NewList()
	admin := GetAdminClient()
	topics := GetTopics(admin)

	ui := UI{
		adminClient:    admin,
		theme:          &theme,
		app:            app,
		sidePane:       sidePane,
		view:           tview.NewGrid(),
		consumersTable: NewSearchableTable(sidePane, app),
		topicsTable:    NewSearchableTable(sidePane, app),
		topics:         topics,
	}

	ui.sidePane.SetBorder(true)
	ui.sidePane.SetBackgroundColor(ui.theme.Background)
	ui.sidePane.SetMainTextColor(ui.theme.Foreground)
	// ui.sidePane.SetSelectedBackgroundColor(ui.theme.PrimaryColor)
	ui.sidePane.SetSelectedTextColor(ui.theme.Foreground)
	ui.sidePane.SetSelectedStyle(tcell.StyleDefault.Attributes(tcell.AttrUnderline))

	ui.view.SetBorder(true)
	ui.view.SetBackgroundColor(ui.theme.Background)

	ui.consumersTable.Table.SetTitle(" Consumer Groups ")
	ui.consumersTable.Table.SetTitleAlign(0)
	ui.consumersTable.Table.SetBorder(true)
	ui.consumersTable.Table.SetBackgroundColor(ui.theme.Background)
	ui.consumersTable.Table.SetSelectable(true, false)
	ui.consumersTable.Table.SetSeparator('┆')

	ui.consumersTable.SearchBox.SetLabel(" > ")
	ui.consumersTable.SearchBox.SetBorder(true)
	ui.consumersTable.SearchBox.SetBorderColor(ui.theme.Foreground)
	ui.consumersTable.SearchBox.SetBackgroundColor(ui.theme.Background)
	ui.consumersTable.SearchBox.SetFieldBackgroundColor(ui.theme.Background)

	ui.topicsTable.Table.SetTitle(" Topics ")
	ui.topicsTable.Table.SetTitleAlign(0)
	ui.topicsTable.Table.SetBorder(true)
	ui.topicsTable.Table.SetBackgroundColor(ui.theme.Background)
	ui.topicsTable.Table.SetSelectable(true, false)
	ui.topicsTable.Table.SetSeparator('┆')

	ui.topicsTable.SearchBox.SetLabel(" > ")
	ui.topicsTable.SearchBox.SetBorder(true)
	ui.topicsTable.SearchBox.SetBorderColor(ui.theme.Foreground)
	ui.topicsTable.SearchBox.SetBackgroundColor(ui.theme.Background)
	ui.topicsTable.SearchBox.SetFieldBackgroundColor(ui.theme.Background)

	ui.sidePane.AddItem("Brokers", "", '1', func() { showBrokersView(&ui) })
	ui.sidePane.AddItem("Topics", "", '2', func() { showTopicsView(&ui) })
	ui.sidePane.AddItem("Consumers", "", '3', func() { showConsumersView(&ui) })

	main := tview.NewFlex()
	main.SetTitle(" Kafka TUI ")
	main.AddItem(ui.sidePane, 20, 0, true)
	main.AddItem(ui.view, 0, 1, true)

	app.SetFocus(ui.sidePane)
	app.EnableMouse(true)

	if err := app.SetRoot(main, true).Run(); err != nil {
		panic(err)
	}
}

func showBrokersView(ui *UI) {

}

func showConsumersView(ui *UI) {
	ui.topics = GetTopics(ui.adminClient)
	ui.consumerGroups = GetConsumersGroups(ui.adminClient)
	ui.consumerGroupsDescriptions = GetConsumersGroupsDescription(ui.adminClient, ui.consumerGroups)

	ui.view.SetBorder(false)

	ui.app.SetFocus(ui.consumersTable.Table)
	ui.view.AddItem(ui.consumersTable.Container, 0, 0, 1, 1, 0, 0, true)

	ui.consumersTable.Table.Clear()

	ui.consumersTable.SetColumnNames([]string{
		" Group ID   ",
		" N° Members   ",
		" N° Partitions   ",
		" Lag   ",
		" Coordinator   ",
		" State   ",
	}, ui.theme.PrimaryColor)

	i := 1
	for group, info := range ui.consumerGroups {

		var description *sarama.GroupDescription
		for _, item := range ui.consumerGroupsDescriptions {
			if item.GroupId == group {
				description = item
				break
			}
		}

		stateCell := tview.NewTableCell(" " + description.State + "   ")
		if description.State == "Stable" {
			stateCell = stateCell.SetTextColor(ui.theme.InfoColor)
		} else {
			stateCell = stateCell.SetTextColor(ui.theme.ErrorColor)
		}

		ui.consumersTable.Table.SetCell(i, 0, tview.NewTableCell(" "+group+"   ").SetTextColor(ui.theme.InEvidenceColor))
		ui.consumersTable.Table.SetCell(i, 1, tview.NewTableCell(" "+strconv.Itoa(len(description.Members))+"   ").SetTextColor(ui.theme.Foreground))
		ui.consumersTable.Table.SetCell(i, 2, tview.NewTableCell(" "+info+"   ").SetTextColor(ui.theme.Foreground))
		ui.consumersTable.Table.SetCell(i, 5, stateCell)
		i++
	}
}

func showTopicsView(ui *UI) {
	topics := GetTopics(ui.adminClient)
	ui.topics = topics
	ui.topicsMetadata = GetTopicsMetadata(ui.adminClient, topics)

	topicNames := []string{}
	for topic := range ui.topics {
		topicNames = append(topicNames, topic)
	}
	ui.topicsSize = GetTopicsSize(ui.adminClient, topicNames)

	ui.view.SetBorder(false)

	ui.app.SetFocus(ui.topicsTable.Table)
	ui.view.AddItem(ui.topicsTable.Container, 0, 0, 1, 1, 0, 0, true)

	ui.topicsTable.Table.Clear()

	ui.topicsTable.SetColumnNames([]string{
		" Name   ",
		" Internal   ",
		" Partitions   ",
		" Out of sync replicas   ",
		" Replication factor   ",
		" Messages   ",
		" Size   ",
	}, ui.theme.PrimaryColor)

	i := 1
	for topic, info := range ui.topics {
		partitions := strconv.FormatInt(int64(info.NumPartitions), 10)
		repFactor := strconv.Itoa(int(info.ReplicationFactor))

		var isInternal string
		for _, item := range ui.topicsMetadata {
			if item.Name == topic {
				isInternal = strconv.FormatBool(item.IsInternal)
				break
			}
		}

		ui.topicsTable.Table.SetCell(i, 0, tview.NewTableCell(" "+topic+"   ").SetTextColor(ui.theme.InEvidenceColor))
		ui.topicsTable.Table.SetCell(i, 1, tview.NewTableCell(" "+isInternal+"   ").SetTextColor(ui.theme.Foreground))
		ui.topicsTable.Table.SetCell(i, 2, tview.NewTableCell(" "+partitions+"   ").SetTextColor(ui.theme.Foreground))
		ui.topicsTable.Table.SetCell(i, 4, tview.NewTableCell(" "+repFactor+"   ").SetTextColor(ui.theme.Foreground))
		ui.topicsTable.Table.SetCell(i, 6, tview.NewTableCell(" "+bytesToString(ui.topicsSize[topic])+"   ").SetTextColor(ui.theme.Foreground))
		i++
	}
}

func bytesToString(bytes int) string {
	m := 1024 //multiplier

	if bytes < 1024 {
		return strconv.Itoa(bytes) + " bytes"
	} else if bytes < m*m {
		return strconv.Itoa(bytes/m) + " KB"
	} else if bytes < m*m*m {
		return strconv.Itoa((bytes/m)/m) + " MB"
	} else if bytes < m*m*m*m {
		return strconv.Itoa(((bytes/m)/m)/m) + " GB"
	} else {
		return strconv.Itoa((((bytes/m)/m)/m)/m) + " TB"
	}
}
