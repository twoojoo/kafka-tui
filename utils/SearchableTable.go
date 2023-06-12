package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
)

type hiddenRow struct {
	index  int
	values []tview.TableCell
}

type SearchableTable struct {
	app              *tview.Application
	Container        *tview.Flex
	Table            *tview.Table
	SearchBox        *tview.InputField
	columnNames      []string
	searchableColumn int
	hiddenRows       []hiddenRow
}

func NewSearchableTable(sidePane *tview.List, app *tview.Application) *SearchableTable {
	t := SearchableTable{
		app:              app,
		Container:        tview.NewFlex(),
		Table:            tview.NewTable(),
		SearchBox:        tview.NewInputField(),
		searchableColumn: 0,
		hiddenRows:       []hiddenRow{},
		columnNames:      []string{},
	}

	t.updateSearchableColumn(t.searchableColumn)

	t.Container.SetDirection(0)
	// t.Container.AddItem(t.SearchBox, 3, 1, false)
	t.Container.AddItem(t.Table, 0, 1, true)

	t.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '\\' {
			t.app.SetFocus(t.SearchBox)
		}

		return event
	})

	t.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		rune := event.Rune()

		if rune == '\\' {
			t.app.SetFocus(t.SearchBox)
		} else if key == tcell.KeyEscape ||
			key == tcell.KeyBackspace ||
			key == tcell.KeyTab {
			app.SetFocus(sidePane)
		}

		return event
	})

	t.SearchBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown {
			t.app.SetFocus(t.Table)
			t.Table.Select(1, 0)
		}

		return event
	})

	t.SearchBox.SetDoneFunc(func(key tcell.Key) {
		t.app.SetFocus(t.Table)
		t.Table.Select(1, 0)
	})

	return &t
}

func (t *SearchableTable) SetColumnNames(columns []string, color tcell.Color) *SearchableTable {
	t.columnNames = columns

	t.printColumnNames(color)

	return t
}

func (t *SearchableTable) printColumnNames(color tcell.Color) *SearchableTable {
	for i := 0; i < len(t.columnNames); i++ {
		t.Table.SetCell(0, i, tview.NewTableCell(t.columnNames[i]).SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold)).SetTextColor(color).SetSelectable(false))
	}

	return t
}

func (t *SearchableTable) SetSearchableColumn(col int) *SearchableTable {
	t.searchableColumn = col
	return t
}

func (t *SearchableTable) updateSearchableColumn(col int) {
	t.SearchBox.SetChangedFunc(func(text string) {
		// colNum := t.Table.GetColumnCount()

		// for _, hr := range t.hiddenRows {
		// 	t.Table.InsertRow(hr.index)

		// 	for i := 0; i < colNum; i++ {
		// 		t.Table.SetCell(hr.index, i, &hr.values[i])
		// 	}
		// }

		// t.hiddenRows = []hiddenRow{}

		//one because of column names
		best, _ := t.Table.GetSelection()
		rows := t.Table.GetRowCount()
		for i := 1; i < rows; i++ {
			rank := fuzzy.RankMatchFold(text, t.Table.GetCell(i, col).Text)
			if rank > best {
				best = rank
			}
		}

		// t.SearchBox.SetText(string(best))
		t.Table.Select(best, 0)
	})
}
