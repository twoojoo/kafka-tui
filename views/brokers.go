package views

import (
	"strconv"

	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/kafka"
	"github.com/twoojoo/ktui/types"
)


func ShowBrokersView(ui *types.UI) {
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
