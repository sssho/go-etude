package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Tui struct {
	app         *tview.Application
	pages       *tview.Pages
	pageTitle   []string
	currentPage int
}

func NewTui() *Tui {
	t := Tui{}
	t.app = tview.NewApplication()
	t.pages = tview.NewPages()
	t.pageTitle = []string{"form1", "form2"}
	form1 := NewForm(t.app, t.pageTitle[0])
	form2 := NewForm(t.app, t.pageTitle[1])
	t.pages.AddPage(t.pageTitle[0], form1.form, true, true)
	t.pages.AddPage(t.pageTitle[1], form2.form, true, false)

	t.SetKeyBind()
	return &t
}

func (t *Tui) SetKeyBind() {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlO:
			t.currentPage = (t.currentPage + 1) % len(t.pageTitle)
			t.pages.SwitchToPage(t.pageTitle[t.currentPage])
		default:
		}
		return event
	})
}

func (tui *Tui) Run() error {
	return tui.app.SetRoot(tui.pages, true).EnableMouse(true).Run()
}

func main() {
	mf := NewTui()
	if err := mf.Run(); err != nil {
		panic(err)
	}
}
