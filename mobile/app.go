package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/L-LYR/pns/mobile/mqtt"
	"github.com/L-LYR/pns/mobile/ui"
	"github.com/L-LYR/pns/proto/pkg/message"
)

func Run() {
	a := app.New()
	client := mqtt.MustNewClient(
		&mqtt.ClientConfig{
			Name:     "pns_target-12345678-12345",
			AppId:    12345,
			Key:      "test_app_name",
			Secret:   "test_app_name",
			DeviceId: "12345678",
			Broker:   "192.168.137.1",
			Port:     "18830",
			Timeout:  60,
		},
	)
	client.Subscribe(NotificationHandler())
	ui.MainWindow(a).ShowAndRun()
}

func NotificationHandler() mqtt.MessageHandler {
	return func(m *message.Message) {
		fyne.CurrentApp().SendNotification(fyne.NewNotification(m.Title, m.Content))
	}
}
