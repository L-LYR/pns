package ui

import "fyne.io/fyne/v2"

type View interface {
	Name() string
	View(w fyne.Window) fyne.CanvasObject
}

var (
	_ View = (*MainView)(nil)
	_ View = (*InfoView)(nil)
)

var (
	_TopWindow fyne.Window = nil
)

func SetTopWindow(w fyne.Window) { _TopWindow = w }
func GetTopWindow() fyne.Window  { return _TopWindow }
