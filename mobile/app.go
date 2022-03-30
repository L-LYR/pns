package mobile

import (
	"fyne.io/fyne/v2/app"
	"github.com/L-LYR/pns/mobile/ui"
)

func Run() {
	a := app.New()
	ui.MainWindow(a).ShowAndRun()
}
