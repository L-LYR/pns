package ui

import (
	"fyne.io/fyne/v2"
)

func switchWindow(v View) {
	parent := GetTopWindow()
	child := fyne.CurrentApp().NewWindow(v.Name())
	SetTopWindow(child)
	child.SetContent(v.View(child))
	child.Show()
	child.SetOnClosed(func() { SetTopWindow(parent) })
}

func NewMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu(
			"Info",
			fyne.NewMenuItem("Device Info", func() {
				switchWindow(NewInfoView())
			}),
		),
	)
}
