package app

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/components"
	"github.com/twoojoo/ktui/kafka"
	"github.com/twoojoo/ktui/types"
	"github.com/twoojoo/ktui/views"
)

const Version = "0.0.1"
const UpdateRateSec = 5

func Run() {

	args := os.Args
	noTopBar := false

	for _, arg := range args {
		if arg == "--version" || arg == "-v" { //VERSION
			fmt.Println("v" + Version)
			os.Exit(0)
		} else if arg == "--help" || arg == "-h" { //HELP
			fmt.Println("\n  ktui v" + Version)
			fmt.Println("  A TUI for Apache Kafka")
			fmt.Println()
			fmt.Println("  -t, --no-top-bar\t hide top bar")
			fmt.Println("  -v, --version\t\t print ktui version")
			fmt.Println("  -h, --help\t\t help")
			fmt.Println()
			os.Exit(0)
		} else if arg == "--no-top-bar" || arg == "-t" { //NO-TOP-BAR
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
		BrokersTable:   components.NewSearchableTable(SidePane, app),
		ConsumersTable: components.NewSearchableTable(SidePane, app),
		TopicsTable:    components.NewSearchableTable(SidePane, app),
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
		views.ShowTopicDetail(&ui, topic)
	})

	ui.SidePane.AddItem("Brokers", "", '1', func() {
		ui.View.Clear()
		ui.View.AddItem(ui.BrokersTable.Container, 0, 2, false)
		ui.UpdateFunc = views.ShowBrokersView
		views.ShowBrokersView(&ui)
		ui.App.SetFocus(ui.BrokersTable.Table)
	})

	ui.SidePane.AddItem("Topics", "", '2', func() {
		ui.View.Clear()
		ui.View.AddItem(ui.TopicsTable.Container, 0, 2, false)

		ui.UpdateFunc = func(ui *types.UI) {
			focused := ui.App.GetFocus()
			defer ui.App.SetFocus(focused)

			views.ShowTopicsView(ui)
			row, _ := ui.TopicsTable.Table.GetSelection()
			topic := ui.TopicsTable.Table.GetCell(row, 0).Text
			views.ShowTopicDetail(ui, topic)
		}

		views.ShowTopicsView(&ui)
		ui.TopicsTable.Table.Select(1, 0)
		topic := ui.TopicsTable.Table.GetCell(1, 0).Text
		views.ShowTopicDetail(&ui, topic)

		ui.SidePane.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			key := event.Key()
			if key == tcell.KeyRight {
				ui.App.SetFocus(ui.TopicsTable.Table)
			} else if key == tcell.KeyLeft {
				ui.App.SetFocus(ui.TopicDetail)
			}

			return event
		})

		ui.App.SetFocus(ui.TopicsTable.Table)
	})

	ui.SidePane.AddItem("Consumers", "", '3', func() {
		ui.View.Clear()
		ui.View.AddItem(ui.ConsumersTable.Container, 0, 2, false)
		ui.UpdateFunc = views.ShowConsumersView
		views.ShowConsumersView(&ui)
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

func getTitle() string {
	title := " ╷  _  _______ _   _ ___"
	title += "\n │ | |/ /_   _| | | |_ _|"
	title += "\n │ | ' /  | | | | | || |"
	title += "\n │ | . \\  | | | |_| || |  v" + Version
	title += "\n │ |_|\\_\\ |_|  \\___/|___| by twoojoo"
	title += "\n └───────────────────────────────────────"
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
