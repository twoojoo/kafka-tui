package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/components"
	"github.com/twoojoo/ktui/kafka"
	"github.com/twoojoo/ktui/types"
)

func Init() *types.UI {
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
		MainContainer:  tview.NewFlex(),
		AdminClient:    admin,
		Theme:          &Theme,
		App:            app,
		SidePane:       SidePane,
		CentralView:    tview.NewFlex(),
		Brokers:        brokers,
		ControllerId:   controllerId,
		BrokersTable:   components.NewSearchableTable(SidePane, app),
		ConsumersTable: components.NewSearchableTable(SidePane, app),
		TopicsTable:    components.NewSearchableTable(SidePane, app),
		Topics:         topics,
		UpdateFunc:     func(*types.UI) {},
	}

	return &ui
}