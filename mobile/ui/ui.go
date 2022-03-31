package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/L-LYR/pns/mobile/push_sdk"
	"github.com/L-LYR/pns/proto/pkg/message"
)

var (
	log = binding.NewStringList()
)

func MainContainer() *fyne.Container {
	pane := widget.NewListWithData(
		log,
		func() fyne.CanvasObject {
			return widget.NewLabel("test")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	paneBox := container.NewVScroll(pane)
	paneBox.SetMinSize(fyne.NewSize(0, 300))
	push_sdk.RegisterGlobalLogHandler(func(s string) {
		log.Append(s)
		pane.ScrollToBottom()
	})
	return container.NewVBox(
		paneBox,
		widget.NewButton("Send Myself a Push", func() {
			push_sdk.PushMyself("Hello", "World")
		}),
		widget.NewButton("Get Token", func() {
			push_sdk.GetToken()
		}),
		widget.NewButton("Subscribe Personal Push", func() {
			push_sdk.RegisterPersonalPushHandler(func(m *message.Message) {
				fyne.CurrentApp().SendNotification(fyne.NewNotification(m.Title, m.Content))
			})
		}),
		widget.NewButton("Subscribe Broadcast Push", func() {
			push_sdk.RegisterBroadcastPushHandler(func(m *message.Message) {
				fyne.CurrentApp().SendNotification(fyne.NewNotification(m.Title, m.Content))
			})
		}),
	)
}
