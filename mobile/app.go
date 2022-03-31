package mobile

import (
	"fyne.io/fyne/v2/app"
	"github.com/L-LYR/pns/mobile/ui"
)

func Run() {
	a := app.New()
	a.SetIcon(ui.ResourceLogoPng)
	window := a.NewWindow("PNS Mobile")
	window.SetContent(
		ui.MainContainer(),
	)
	window.ShowAndRun()
}
