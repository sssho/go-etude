package main

import (
	"fmt"
	"os/exec"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func runSubprocess() string {
	cmd := exec.Command("sh", "-c", "pwd | fzf-tmux")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

type FormParam struct {
	Title    string
	FName    string
	LName    string
	Age18    bool
	Password string
}

type FormLabel struct {
	title    string
	fName    string
	lName    string
	age18    string
	password string
	save     string
	quit     string
}

type MyForm struct {
	app   *tview.Application
	form  *tview.Form
	label FormLabel
	Param FormParam
}

func (f MyForm) GetTextData(label string) string {
	item := f.form.GetFormItemByLabel(label)
	var data string
	switch t := item.(type) {
	case *tview.DropDown:
		_, data = t.GetCurrentOption()
	case *tview.InputField:
		data = t.GetText()
	default:
		data = "default"
	}
	return data
}

func (f MyForm) IsChecked(label string) bool {
	item := f.form.GetFormItemByLabel(label)
	var isChecked bool
	switch t := item.(type) {
	case *tview.Checkbox:
		isChecked = t.IsChecked()
	default:
		isChecked = false
	}
	return isChecked
}

func (f MyForm) SetSaveFunc(handler func()) {
	saveButton := f.form.GetButton(0)
	saveButton.SetSelectedFunc(handler)
}

func (f *MyForm) SetFormParam() {
	f.Param.Title = f.GetTextData(f.label.title)
	f.Param.FName = f.GetTextData(f.label.fName)
	f.Param.LName = f.GetTextData(f.label.lName)
	f.Param.Age18 = f.IsChecked(f.label.age18)
	f.Param.Password = f.GetTextData(f.label.password)
}

func (f *MyForm) SetKeyBind() {
	f.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			i, _ := f.form.GetFocusedItemIndex()
			item := f.form.GetFormItem(i)
			if xi, ok := item.(*tview.InputField); ok {
				data := runSubprocess()
				xi.SetText(data)
			}
		case tcell.KeyCtrlS:
			_ = f.app.Suspend(func() {
				fmt.Printf("%v\n", f.Param)
			})
		default:
			// _ = app.Suspend(runSubprocess)
		}
		return event
	})
}

func NewForm(app *tview.Application, title string) *MyForm {
	f := MyForm{}
	f.app = app
	f.label = FormLabel{
		title:    "Title",
		fName:    "First name",
		lName:    "Last name",
		age18:    "Age 18+",
		password: "Password",
		save:     "Save",
		quit:     "Quit",
	}
	f.Param = FormParam{}
	f.form = tview.NewForm().
		AddDropDown(f.label.title, []string{"Mr.", "Ms."}, 0, nil).
		AddInputField(f.label.fName, "", 20, nil, nil).
		AddInputField(f.label.lName, "", 20, nil, nil).
		AddCheckbox(f.label.age18, false, nil).
		AddPasswordField(f.label.password, "", 10, '*', nil).
		AddButton(f.label.save, nil).
		AddButton(f.label.quit, nil)
	f.form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	f.SetKeyBind()

	f.SetSaveFunc(func() {
		(&f).SetFormParam()
	})

	return &f
}
