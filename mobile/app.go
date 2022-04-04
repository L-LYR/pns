package mobile

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/L-LYR/pns/mobile/ui"
)

func Run() {
	a := app.New()
	a.SetIcon(ui.ResourceLogoPng)
	a.Settings().SetTheme(theme.LightTheme())

	window := a.NewWindow("PNS Mobile")
	window.SetMainMenu(ui.NewMenu())
	window.SetMaster()
	ui.SetTopWindow(window)
	window.SetContent(ui.NewMainView().View(window))
	window.ShowAndRun()
}
