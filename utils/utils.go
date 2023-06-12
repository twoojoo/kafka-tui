package utils

import (
	"sort"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/types"
)

func BytesToString(bytes int) string {
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

func SortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func BuildDetailText(ui *types.UI, text string) *tview.TextView {
	box := tview.NewTextView()
	box.SetText(text)
	box.SetBackgroundColor(ui.Theme.Background)
	box.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))
	return box
}