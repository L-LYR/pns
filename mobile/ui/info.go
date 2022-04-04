package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/L-LYR/pns/mobile/push_sdk"
)

type InfoView struct{}

func NewInfoView() *InfoView { return &InfoView{} }

func (i *InfoView) Name() string { return "Target Info" }
func (i *InfoView) View(w fyne.Window) fyne.CanvasObject {
	info := widget.NewTextGrid()
	info.SetText(
		"Device ID:\n" +
			"  " + push_sdk.GetConfig().DeviceId + "\n" +
			"App ID:\n" +
			"  " + strconv.FormatInt(int64(push_sdk.GetConfig().AppId), 10) + "\n",
	)
	return container.NewBorder(
		nil, nil, nil, nil,
		info,
	)
}
