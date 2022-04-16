package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/L-LYR/pns/mobile/push_sdk"
	"github.com/L-LYR/pns/proto/pkg/message"
)

type MainView struct{}

func NewMainView() *MainView { return &MainView{} }

func (m *MainView) Name() string { return "PNS Mobile" }
func (m *MainView) View(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("PNS Mobile")
	title.Alignment = fyne.TextAlignCenter

	logPane := widget.NewListWithData(
		log,
		func() fyne.CanvasObject {
			label := widget.NewLabel("empty")
			label.Wrapping = fyne.TextTruncate
			return label
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)

	log.AddListener(binding.NewDataListener(
		func() { logPane.ScrollToBottom() },
	))

	pushButton := widget.NewButton("Send Myself a Push", func() {
		PushMyself("Hello", "World")
	})

	targetInfoButtons := container.NewGridWithColumns(
		2,
		widget.NewButton("Update Info", func() {
			push_sdk.UpdateTargetInfo()
		}),
		widget.NewButton("Get Token", func() {
			push_sdk.GetToken()
		}),
	)

	subscribeButtons := container.NewGridWithColumns(
		2,
		widget.NewButton("Subscribe PPush", func() {
			push_sdk.RegisterPersonalPushHandler(func(m *message.Message) error {
				fyne.CurrentApp().SendNotification(fyne.NewNotification(m.Title, m.Content))
				return nil
			})
		}),
		widget.NewButton("Subscribe BPush", func() {
			push_sdk.RegisterBroadcastPushHandler(func(m *message.Message) error {
				fyne.CurrentApp().SendNotification(fyne.NewNotification(m.Title, m.Content))
				return nil
			})
		}),
	)

	themeSelector := container.NewGridWithColumns(
		2,
		widget.NewButton("Dark", func() {
			fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme()) //lint:ignore SA1019 no reason
		}),
		widget.NewButton("Light", func() {
			fyne.CurrentApp().Settings().SetTheme(theme.LightTheme()) //lint:ignore SA1019 no reason
		}),
	)

	return container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
		),
		container.NewGridWithRows(
			4,
			targetInfoButtons,
			pushButton,
			subscribeButtons,
			themeSelector,
		),
		nil,
		nil,
		container.NewVScroll(logPane),
	)
}
