package ui

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"github.com/L-LYR/pns/mobile/push_sdk"
	"github.com/L-LYR/pns/mobile/push_sdk/net/http"
)

var (
	log = binding.NewStringList()

	LogHandler = func(s string, v ...interface{}) {
		log.Append(fmt.Sprintf(s, v...))
	}

	bizHTTPClient = http.MustNewHTTPClient("http://192.168.137.1:10087")
)

func PushMyself(title string, content string) {
	payload, err := bizHTTPClient.POST("/push/direct", http.Payload{
		"deviceId": push_sdk.GetConfig().GetDeviceId(),
		"appId":    push_sdk.GetConfig().GetAppId(),
		"title":    title,
		"content":  content,
	})
	if err != nil {
		LogHandler("Error: %s", err.Error())
	} else {
		if s, ok := payload["pushTaskId"]; ok {
			LogHandler("Info: Task ID: %+v", s)
		} else {
			LogHandler("Error: no push task id")
		}
	}
}

func Run() {
	a := app.NewWithID("PNS")
	a.SetIcon(ResourceLogoPng)
	a.Settings().SetTheme(theme.LightTheme())
	// Init Push SDK
	window := a.NewWindow("PNS Mobile")
	window.SetMainMenu(NewMenu())
	window.SetMaster()
	SetTopWindow(window)
	window.SetContent(NewMainView().View(window))
	window.ShowAndRun()
}
