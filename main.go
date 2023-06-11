package main

import (
	"os"
	"sort"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	//graphic stuff
	theme          *Theme
	app            *tview.Application
	sidePane       *tview.List
	view           *tview.Grid
	brokersTable   *SearchableTable
	consumersTable *SearchableTable
	topicsTable    *SearchableTable

	//kafka stuff
	adminClient                *sarama.ClusterAdmin
	brokers                    []*sarama.Broker
	controllerId               int32
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

const Version = "0.0.1"

func main() {

	args := os.Args
	noTopBar := false

	for _, arg := range args {
		if arg == "--no-top-bar" {
			noTopBar = true
		}
	}

	app := tview.NewApplication()

	theme := Theme{
		Background:      tcell.ColorReset,
		Foreground:      tcell.ColorWhite,
		PrimaryColor:    tcell.ColorHotPink,
		InEvidenceColor: tcell.ColorYellow,
		InfoColor:       tcell.ColorGreen,
		ErrorColor:      tcell.ColorRed,
		WarningColor:    tcell.ColorYellow,
	}

	sidePane := tview.NewList()
	admin := GetAdminClient()
	topics := GetTopics(admin)
	brokers, controllerId := GetBrokers(admin)

	ui := UI{
		adminClient:    admin,
		theme:          &theme,
		app:            app,
		sidePane:       sidePane,
		view:           tview.NewGrid(),
		brokers:        brokers,
		controllerId:   controllerId,
		brokersTable:   NewSearchableTable(sidePane, app),
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

	ui.brokersTable.Table.SetTitle(" Brokers ")
	ui.brokersTable.Table.SetTitleAlign(0)
	ui.brokersTable.Table.SetBorder(true)
	ui.brokersTable.Table.SetBackgroundColor(ui.theme.Background)
	ui.brokersTable.Table.SetSelectable(true, false)
	ui.brokersTable.Table.SetSeparator('┆')

	ui.brokersTable.SearchBox.SetLabel(" > ")
	ui.brokersTable.SearchBox.SetBorder(true)
	ui.brokersTable.SearchBox.SetBorderColor(ui.theme.Foreground)
	ui.brokersTable.SearchBox.SetBackgroundColor(ui.theme.Background)
	ui.brokersTable.SearchBox.SetFieldBackgroundColor(ui.theme.Background)

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

	main2 := tview.NewFlex()
	main2.SetTitle(" Kafka TUI ")
	main2.AddItem(ui.sidePane, 20, 0, true)
	main2.AddItem(ui.view, 0, 1, false)

	main1 := tview.NewFlex()
	main1.SetDirection(0)

	if !noTopBar {
		topBar := tview.NewFlex()
		topBar.SetBorder(true)
		topBar.SetBackgroundColor(ui.theme.Background)

		titleBar := tview.NewTextView()
		titleBar.SetText(getTitle())
		titleBar.SetBackgroundColor(ui.theme.Background)
		titleBar.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
		titleBar.SetTextColor(ui.theme.PrimaryColor)

		hotkeysText := tview.NewTextView()
		hotkeysText.SetText(getHotkeysText())
		hotkeysText.SetTextAlign(2)
		hotkeysText.SetBackgroundColor(ui.theme.Background)
		hotkeysText.SetTextColor(ui.theme.Foreground)
		hotkeysText.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))

		hotkeysKeys := tview.NewTextView()
		hotkeysKeys.SetText(getHotkeys())
		hotkeysKeys.SetTextAlign(2)
		hotkeysKeys.SetBackgroundColor(ui.theme.Background)
		hotkeysKeys.SetTextColor(ui.theme.InEvidenceColor)

		topBar.AddItem(titleBar, 0, 1, false)
		topBar.AddItem(hotkeysText, 17, 0, false)
		topBar.AddItem(hotkeysKeys, 10, 0, false)

		main1.AddItem(topBar, 8, 0, false)
	}

	main1.AddItem(main2, 0, 1, true)

	app.SetFocus(ui.sidePane)
	app.EnableMouse(true)

	if err := app.SetRoot(main1, true).Run(); err != nil {
		panic(err)
	}
}

func showBrokersView(ui *UI) {
	brokers, controllerId := GetBrokers(ui.adminClient)
	ui.brokers = brokers
	ui.controllerId = controllerId

	ui.view.SetBorder(false)

	ui.app.SetFocus(ui.brokersTable.Table)
	ui.view.AddItem(ui.brokersTable.Container, 0, 0, 1, 1, 0, 0, true)

	ui.brokersTable.Table.Clear()

	for i, broker := range ui.brokers {

		ui.brokersTable.SetColumnNames([]string{
			" ID   ",
			" Address   ",
			// " N° Partitions   ",
			// " Lag   ",
			// " Coordinator   ",
			// " State   ",
		}, ui.theme.PrimaryColor)

		// i := 1
		// for group, info := range ui.consumerGroups {

		// 	var description *sarama.GroupDescription
		// 	for _, item := range ui.consumerGroupsDescriptions {
		// 		if item.GroupId == group {
		// 			description = item
		// 			break
		// 		}
		// 	}

		// 	stateCell := tview.NewTableCell(" " + description.State + "   ")
		// 	if description.State == "Stable" {
		// 		stateCell = stateCell.SetTextColor(ui.theme.InfoColor)
		// 	} else {
		// 		stateCell = stateCell.SetTextColor(ui.theme.ErrorColor)
		// 	}

		ui.brokersTable.Table.SetCell(i+1, 0, tview.NewTableCell(" "+strconv.Itoa(int(broker.ID()))+"   ").SetTextColor(ui.theme.InEvidenceColor))
		ui.brokersTable.Table.SetCell(i+1, 1, tview.NewTableCell(" "+broker.Addr()+"   ").SetTextColor(ui.theme.Foreground))
		// 	ui.consumersTable.Table.SetCell(i, 2, tview.NewTableCell(" "+info+"   ").SetTextColor(ui.theme.Foreground))
		// 	ui.consumersTable.Table.SetCell(i, 5, stateCell)
		// 	i++
	}
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

	sort.Strings(topicNames)
	ui.topicsSize = GetTopicsSize(ui.adminClient, topicNames)

	ui.view.SetBorder(false)

	ui.app.SetFocus(ui.topicsTable.Table)
	ui.view.AddItem(ui.topicsTable.Container, 0, 0, 1, 1, 0, 0, true)

	// ui.topicsTable.Table.Clear()

	ui.topicsTable.SetColumnNames([]string{
		" Name   ",
		" Internal   ",
		" Partitions   ",
		// " Out of sync replicas   ",
		" Replication factor   ",
		" Messages   ",
		" Size   ",
	}, ui.theme.PrimaryColor)

	i := 1
	for _, topic := range topicNames {
		info := ui.topics[topic]
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
		ui.topicsTable.Table.SetCell(i, 3, tview.NewTableCell(" "+repFactor+"   ").SetTextColor(ui.theme.Foreground))
		ui.topicsTable.Table.SetCell(i, 5, tview.NewTableCell(" "+bytesToString(ui.topicsSize[topic])+"   ").SetTextColor(ui.theme.Foreground))
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

func getTitle() string {
	title := "  _  _______ _   _ ___"
	title += "\n | |/ /_   _| | | |_ _|"
	title += "\n | ' /  | | | | | || |"
	title += "\n | . \\  | | | |_| || |"
	title += "\n |_|\\_\\ |_|  \\___/|___| v" + Version + " (by twoojoo)"
	return title
}

func getHotkeysText() string {
	htkTxt := ""
	htkTxt += "\nmove through tabs"
	htkTxt += "\nfocus search bar"
	htkTxt += "\nselect element"
	htkTxt += "\nmove"

	return htkTxt
}

func getHotkeys() string {
	htkTxt := ""
	htkTxt += "\nTab  "
	htkTxt += "\n\\  "
	htkTxt += "\nEnter  "
	htkTxt += "\n🡡 🡣  "

	return htkTxt
}
