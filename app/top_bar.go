package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/types"
)

func AddTopBar(ui *types.UI) {
	topBar := tview.NewFlex()
	topBar.SetBorder(true)
	topBar.SetBackgroundColor(ui.Theme.Background)

	titleBar := tview.NewTextView()
	titleBar.SetText(GetTitle())
	titleBar.SetBackgroundColor(ui.Theme.Background)
	titleBar.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
	titleBar.SetTextColor(ui.Theme.PrimaryColor)

	htkTxt1 := tview.NewTextView()
	htkTxt1.SetText(GetHotkeysText1())
	htkTxt1.SetTextAlign(2)
	htkTxt1.SetBackgroundColor(ui.Theme.Background)
	htkTxt1.SetTextColor(ui.Theme.Foreground)
	htkTxt1.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))

	ktkKeys1 := tview.NewTextView()
	ktkKeys1.SetText(GetHotkeysKeys1())
	ktkKeys1.SetTextAlign(2)
	ktkKeys1.SetBackgroundColor(ui.Theme.Background)
	ktkKeys1.SetTextColor(ui.Theme.InEvidenceColor)

	htkTxt2 := tview.NewTextView()
	htkTxt2.SetText(GetHotkeysText2())
	htkTxt2.SetTextAlign(2)
	htkTxt2.SetBackgroundColor(ui.Theme.Background)
	htkTxt2.SetTextColor(ui.Theme.Foreground)
	htkTxt2.SetTextStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))

	ktkKeys2 := tview.NewTextView()
	ktkKeys2.SetText(GetHotkeysKeys2())
	ktkKeys2.SetTextAlign(2)
	ktkKeys2.SetBackgroundColor(ui.Theme.Background)
	ktkKeys2.SetTextColor(ui.Theme.InEvidenceColor)

	topBar.AddItem(titleBar, 0, 1, false)
	topBar.AddItem(htkTxt1, 20, 0, false)
	topBar.AddItem(ktkKeys1, 12, 0, false)
	topBar.AddItem(htkTxt2, 23, 0, false)
	topBar.AddItem(ktkKeys2, 7, 0, false)

	ui.MainContainer.AddItem(topBar, 8, 0, false)
}
