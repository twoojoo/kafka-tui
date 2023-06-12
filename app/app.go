package app

import (
	"os"
	"sort"
	"strconv"
	"strings"

	// "strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/kafka"
	"github.com/twoojoo/ktui/utils"
	"github.com/twoojoo/ktui/types"
	"github.com/twoojoo/ktui/views"
)

const Version = "0.0.1"
const UpdateRateSec = 5

func Run() {

	args := os.Args
	noTopBar := false

	for _, arg := range args {
		if arg == "--no-top-bar" {
			noTopBar = true
		}
	}

	app := tview.NewApplication()

	Theme := types.Theme{
		Background:      tcell.ColorReset,
		Foreground:      tcell.ColorWhite,
		PrimaryColor:    tcell.ColorHotPink,
		InEvidenceColor: tcell.ColorYellow,
		InfoColor:       tcell.ColorGreen,
		ErrorColor:      tcell.ColorRed,
		WarningColor:    tcell.ColorYellow,
	}

	SidePane := tview.NewList()
	admin := kafka.GetAdminClient()
	topics := kafka.GetTopics(admin)
	brokers, controllerId := kafka.GetBrokers(admin)

	ui := types.UI{
		AdminClient:    admin,
		Theme:          &Theme,
		App:            app,
		SidePane:       SidePane,
		View:           tview.NewFlex(),
		Brokers:        brokers,
		ControllerId:   controllerId,
		BrokersTable:   utils.NewSearchableTable(SidePane, app),
		ConsumersTable: utils.NewSearchableTable(SidePane, app),
		TopicsTable:    utils.NewSearchableTable(SidePane, app),
		Topics:         topics,
		UpdateFunc:     func(*types.UI) {},
	}

	//set the refresh goroutine
	go func(ui *types.UI) {
		for {
			time.Sleep(UpdateRateSec * time.Second)
			ui.UpdateFunc(ui)
			ui.App.Draw()
		}
	}(&ui)

	ui.SidePane.SetBorder(true)
	ui.SidePane.SetBackgroundColor(ui.Theme.Background)
	ui.SidePane.SetMainTextColor(ui.Theme.Foreground)
	ui.SidePane.SetSelectedTextColor(ui.Theme.Foreground)
	ui.SidePane.SetSelectedStyle(tcell.StyleDefault.Attributes(tcell.AttrUnderline))

	ui.View.SetBorder(true)
	ui.View.SetBackgroundColor(ui.Theme.Background)

	ui.BrokersTable.Table.SetTitle(" Brokers ")
	ui.BrokersTable.Table.SetTitleAlign(0)
	ui.BrokersTable.Table.SetBorder(true)
	ui.BrokersTable.Table.SetBackgroundColor(ui.Theme.Background)
	ui.BrokersTable.Table.SetSelectable(true, false)
	ui.BrokersTable.Table.SetSeparator('┆')

	ui.BrokersTable.SearchBox.SetLabel(" > ")
	ui.BrokersTable.SearchBox.SetBorder(true)
	ui.BrokersTable.SearchBox.SetBorderColor(ui.Theme.Foreground)
	ui.BrokersTable.SearchBox.SetBackgroundColor(ui.Theme.Background)
	ui.BrokersTable.SearchBox.SetFieldBackgroundColor(ui.Theme.Background)

	ui.ConsumersTable.Table.SetTitle(" Consumer Groups ")
	ui.ConsumersTable.Table.SetTitleAlign(0)
	ui.ConsumersTable.Table.SetBorder(true)
	ui.ConsumersTable.Table.SetBackgroundColor(ui.Theme.Background)
	ui.ConsumersTable.Table.SetSelectable(true, false)
	ui.ConsumersTable.Table.SetSeparator('┆')

	ui.ConsumersTable.SearchBox.SetLabel(" > ")
	ui.ConsumersTable.SearchBox.SetBorder(true)
	ui.ConsumersTable.SearchBox.SetBorderColor(ui.Theme.Foreground)
	ui.ConsumersTable.SearchBox.SetBackgroundColor(ui.Theme.Background)
	ui.ConsumersTable.SearchBox.SetFieldBackgroundColor(ui.Theme.Background)

	ui.TopicsTable.Table.SetTitle(" Topics ")
	ui.TopicsTable.Table.SetTitleAlign(0)
	ui.TopicsTable.Table.SetBorder(true)
	ui.TopicsTable.Table.SetBackgroundColor(ui.Theme.Background)
	ui.TopicsTable.Table.SetSelectable(true, false)
	ui.TopicsTable.Table.SetSeparator('┆')

	ui.TopicsTable.SearchBox.SetLabel(" > ")
	ui.TopicsTable.SearchBox.SetBorder(true)
	ui.TopicsTable.SearchBox.SetBorderColor(ui.Theme.Foreground)
	ui.TopicsTable.SearchBox.SetBackgroundColor(ui.Theme.Background)
	ui.TopicsTable.SearchBox.SetFieldBackgroundColor(ui.Theme.Background)

	ui.TopicsTable.Table.SetSelectionChangedFunc(func(row int, _ int) {
		if row < 1 {
			return
		}

		topic := ui.TopicsTable.Table.GetCell(row, 0).Text
		showTopicDetail(&ui, topic)
	})

	ui.SidePane.AddItem("Brokers", "", '1', func() {
		ui.View.Clear()
		ui.View.AddItem(ui.BrokersTable.Container, 0, 2, false)
		ui.UpdateFunc = showBrokersView
		showBrokersView(&ui)
		ui.App.SetFocus(ui.BrokersTable.Table)
	})

	ui.SidePane.AddItem("Topics", "", '2', func() {
		ui.View.Clear()
		ui.View.AddItem(ui.TopicsTable.Container, 0, 2, false)

		ui.UpdateFunc = func(ui *types.UI) {
			views.ShowTopicsView(ui)
			ui.App.SetFocus(ui.TopicsTable.Table)
			row, _ := ui.TopicsTable.Table.GetSelection()
			topic := ui.TopicsTable.Table.GetCell(row, 0).Text
			showTopicDetail(ui, topic)
		}

		views.ShowTopicsView(&ui)
		ui.App.SetFocus(ui.TopicsTable.Table)
		ui.TopicsTable.Table.Select(1, 0)
		topic := ui.TopicsTable.Table.GetCell(1, 0).Text
		showTopicDetail(&ui, topic)
	})

	ui.SidePane.AddItem("Consumers", "", '3', func() {
		ui.View.Clear()
		ui.View.AddItem(ui.ConsumersTable.Container, 0, 2, false)
		ui.UpdateFunc = showConsumersView
		showConsumersView(&ui)
		ui.App.SetFocus(ui.ConsumersTable.Table)
	})

	ui.SidePane.AddItem("ACLs", "", '4', func() {
		return
	})

	main2 := tview.NewFlex()
	main2.SetTitle(" Kafka TUI ")
	main2.AddItem(ui.SidePane, 20, 0, true)
	main2.AddItem(ui.View, 0, 1, false)

	main1 := tview.NewFlex()
	main1.SetDirection(0)

	if !noTopBar {
		topBar := tview.NewFlex()
		topBar.SetBorder(true)
		topBar.SetBackgroundColor(ui.Theme.Background)

		titleBar := tview.NewTextView()
		titleBar.SetText(getTitle())
		titleBar.SetBackgroundColor(ui.Theme.Background)
		titleBar.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
		titleBar.SetTextColor(ui.Theme.PrimaryColor)
		
		hotkeysText := tview.NewTextView()
		hotkeysText.SetText(getHotkeysText())
		hotkeysText.SetTextAlign(2)
		hotkeysText.SetBackgroundColor(ui.Theme.Background)
		hotkeysText.SetTextColor(ui.Theme.Foreground)
		hotkeysText.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))

		hotkeysKeys := tview.NewTextView()
		hotkeysKeys.SetText(getHotkeys())
		hotkeysKeys.SetTextAlign(2)
		hotkeysKeys.SetBackgroundColor(ui.Theme.Background)
		hotkeysKeys.SetTextColor(ui.Theme.InEvidenceColor)

		topBar.AddItem(titleBar, 0, 1, false)
		topBar.AddItem(hotkeysText, 17, 0, false)
		topBar.AddItem(hotkeysKeys, 10, 0, false)

		main1.AddItem(topBar, 8, 0, false)
	}

	main1.AddItem(main2, 0, 1, true)

	app.SetFocus(ui.SidePane)
	app.EnableMouse(true)

	if err := app.SetRoot(main1, true).Run(); err != nil {
		panic(err)
	}
}

func showBrokersView(ui *types.UI) {
	brokers, controllerId := kafka.GetBrokers(ui.AdminClient)
	ui.Brokers = brokers
	ui.ControllerId = controllerId

	ui.View.SetBorder(false)

	ui.BrokersTable.Table.Clear()

	for i, broker := range ui.Brokers {

		ui.BrokersTable.SetColumnNames([]string{
			" ID   ",
			" Address   ",
			// " N° Partitions   ",
			// " Lag   ",
			// " Coordinator   ",
			// " State   ",
		}, ui.Theme.PrimaryColor)

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
		// 		stateCell = stateCell.SetTextColor(ui.Theme.InfoColor)
		// 	} else {
		// 		stateCell = stateCell.SetTextColor(ui.Theme.ErrorColor)
		// 	}

		ui.BrokersTable.Table.SetCell(i+1, 0, tview.NewTableCell(" "+strconv.Itoa(int(broker.ID()))+"   ").SetTextColor(ui.Theme.InEvidenceColor))
		ui.BrokersTable.Table.SetCell(i+1, 1, tview.NewTableCell(" "+broker.Addr()+"   ").SetTextColor(ui.Theme.Foreground))
		// 	ui.ConsumersTable.Table.SetCell(i, 2, tview.NewTableCell(" "+info+"   ").SetTextColor(ui.Theme.Foreground))
		// 	ui.ConsumersTable.Table.SetCell(i, 5, stateCell)
		// 	i++
	}
}

func showConsumersView(ui *types.UI) {
	ui.Topics = kafka.GetTopics(ui.AdminClient)
	ui.ConsumerGroups = kafka.GetConsumersGroups(ui.AdminClient)
	ui.ConsumerGroupsDescriptions = kafka.GetConsumersGroupsDescription(ui.AdminClient, ui.ConsumerGroups)

	ui.View.SetBorder(false)

	ui.ConsumersTable.Table.Clear()

	ui.ConsumersTable.SetColumnNames([]string{
		" Group ID   ",
		" N° Members   ",
		" N° Partitions   ",
		" Lag   ",
		" Coordinator   ",
		" State   ",
	}, ui.Theme.PrimaryColor)

	i := 1
	for group, info := range ui.ConsumerGroups {

		var description *sarama.GroupDescription
		for _, item := range ui.ConsumerGroupsDescriptions {
			if item.GroupId == group {
				description = item
				break
			}
		}

		stateCell := tview.NewTableCell(" " + description.State + "   ")
		if description.State == "Stable" {
			stateCell = stateCell.SetTextColor(ui.Theme.InfoColor)
		} else {
			stateCell = stateCell.SetTextColor(ui.Theme.ErrorColor)
		}

		ui.ConsumersTable.Table.SetCell(i, 0, tview.NewTableCell(" "+group+"   ").SetTextColor(ui.Theme.InEvidenceColor))
		ui.ConsumersTable.Table.SetCell(i, 1, tview.NewTableCell(" "+strconv.Itoa(len(description.Members))+"   ").SetTextColor(ui.Theme.Foreground))
		ui.ConsumersTable.Table.SetCell(i, 2, tview.NewTableCell(" "+info+"   ").SetTextColor(ui.Theme.Foreground))
		ui.ConsumersTable.Table.SetCell(i, 5, stateCell)
		i++
	}
}

func showTopicDetail(ui *types.UI, topic string) {
	topic = strings.Trim(topic, " ")
	info := ui.Topics[topic]
	metadata := getTopicMetadata(ui, topic)

	// ui.TopicDetail.Clear()
	ui.TopicDetail = tview.NewFlex()
	ui.TopicDetail.SetBorder(true)
	ui.TopicDetail.SetBackgroundColor(ui.Theme.Background)
	ui.TopicDetail.SetDirection(0)
	ui.TopicDetail.SetTitle(" Topic detail: ")
	ui.TopicDetail.SetTitleAlign(0)

	ui.View.Clear()
	ui.View.AddItem(ui.TopicsTable.Container, 0, 2, true)
	ui.View.AddItem(ui.TopicDetail, 60, 1, false)

	detailTitle := tview.NewTextView()
	detailTitle.SetText("\n " + topic)
	detailTitle.SetBackgroundColor(ui.Theme.Background)
	detailTitle.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
	detailTitle.SetTextColor(ui.Theme.PrimaryColor)

	kindText := buildDetailText(ui, " - internal: "+strconv.FormatBool(metadata.IsInternal))
	sizeText := buildDetailText(ui, " - size: "+utils.BytesToString(ui.TopicsSize[topic]))
	partitionsText := buildDetailText(ui, " - partitions: "+strconv.Itoa(len(metadata.Partitions)))
	messagesText := buildDetailText(ui, " - messages: ")
	replicaText := buildDetailText(ui, " - rep. factor: "+strconv.Itoa(int(info.ReplicationFactor)))

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
	ui.TopicDetail.AddItem(buildDetailText(ui, ""), 1, 0, false)

	sortedConfigNames := sortMapKeys(info.ConfigEntries)

	cfgSubtitle := tview.NewTextView()
	cfgSubtitle.SetBackgroundColor(ui.Theme.Background)
	cfgSubtitle.SetText(" Config:")
	cfgSubtitle.SetTextColor(ui.Theme.Foreground)
	ui.TopicDetail.AddItem(cfgSubtitle, 1, 0, false)

	for _, name := range sortedConfigNames {
		text := buildDetailText(ui, " - "+name+": "+*info.ConfigEntries[name])
		ui.TopicDetail.AddItem(text, 1, 0, false)
	}

	ui.TopicDetail.AddItem(filler, 0, 1, false)

	detailMenu := tview.NewList()
	detailMenu.AddItem("Edit Config", "", '1', func() {})
	detailMenu.AddItem("Clear Messages", "", '2', func() {})
	detailMenu.AddItem("Recreate Topic", "", '3', func() {})
	detailMenu.AddItem("Remove Topic", "", '4', func() {})
	detailMenu.SetMainTextColor(ui.Theme.Foreground)
	detailMenu.SetBackgroundColor(ui.Theme.Background)
	detailMenu.SetSelectedTextColor(ui.Theme.Foreground)
	detailMenu.SetSelectedStyle(tcell.StyleDefault.Attributes(tcell.AttrUnderline))

	ui.TopicDetail.AddItem(detailMenu, 8, 1, false)
}

func getTitle() string {
	title := " ╷  _  _______ _   _ ___"
	title += "\n │ | |/ /_   _| | | |_ _|"
	title += "\n │ | ' /  | | | | | || |"
	title += "\n │ | . \\  | | | |_| || |"
	title += "\n │ |_|\\_\\ |_|  \\___/|___| v" + Version + " (by twoojoo)"
	title += "\n └────────────────────────────────────────────"
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
	htkTxt += "\nTab   "
	htkTxt += "\n\\   "
	htkTxt += "\nEnter   "
	htkTxt += "\n🡡 🡣   "

	return htkTxt
}

func getTopicMetadata(ui *types.UI, topic string) *sarama.TopicMetadata {
	for _, item := range ui.TopicsMetadata {
		if item.Name == topic {
			return item
		}
	}

	panic("no such topic: " + topic)
}

func buildDetailText(ui *types.UI, text string) *tview.TextView {
	// if fg == nil {
	// 	fg = ui.Theme.Foreground
	// }
	box := tview.NewTextView()
	box.SetText(text)
	box.SetBackgroundColor(ui.Theme.Background)
	box.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))
	return box
}

func getKafkaLogo() string {
	logo := "\n\n"
	logo += "                    @@@@@@\n"
	logo += "                   @@    @@@\n"
	logo += "                   @@    @@@\n"
	logo += "                    @@@@@@     @@@@@@@@\n"
	logo += "                      @@      @@@    @@@\n"
	logo += "                   @@@@@@@@  %@@@@  @@@,\n"
	logo += "                 #@@@    .@@@%   &@@@\n"
	logo += "                 @@@       @@&\n"
	logo += "                  @@@    @@@@@   @@@@\n"
	logo += "                   @@@@@@@   @@@%   @@@\n"
	logo += "                      @@      @@@    @@@\n"
	logo += "                    @@@@@@*    %@@@@@@@\n"
	logo += "                   @@    @@@\n"
	logo += "                   @@    @@@\n"
	logo += "                    @@@@@@"
	return logo
}

func sortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
