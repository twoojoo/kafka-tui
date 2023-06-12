package views

import (
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/kafka"
	"github.com/twoojoo/ktui/types"
)


func ShowConsumersView(ui *types.UI) {
	ui.Topics = kafka.GetTopics(ui.AdminClient)
	ui.ConsumerGroups = kafka.GetConsumersGroups(ui.AdminClient)
	ui.ConsumerGroupsDescriptions = kafka.GetConsumersGroupsDescription(ui.AdminClient, ui.ConsumerGroups)

	ui.CentralView.SetBorder(false)

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