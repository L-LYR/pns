package ui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Hello")
	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
		widget.NewButton("push", TestHandle),
	))
	return w
}

func TestHandle() {
	url := "http://192.168.137.1:10087/push"

	payload := strings.NewReader("{\n\t\"deviceId\": \"12345678\",\n\t\"appId\": 12345,\n\t\"title\": \"Hello\",\n\t\"content\": \"world\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
