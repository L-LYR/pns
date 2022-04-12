package frontend

import (
	"github.com/L-LYR/pns/internal/admin/frontend/components"
	"github.com/L-LYR/pns/internal/admin/frontend/settings"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

var (
	Admin = &app.Handler{
		Author:          "HammerLi",
		Name:            "PNS Admin",
		Title:           "PNS Admin",
		BackgroundColor: settings.BackgroundColor,
		ThemeColor:      settings.BackgroundColor,
		Icon: app.Icon{
			Default: "/web/logo.svg",
		},
		Styles: []string{
			"https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500&display=swap",
			"/web/index.css",
		},
	}
)

func MustRegisterFrontendRouters() {
	settings.MustLoadConfig()
	ui.BaseHPadding = settings.BaseHPadding
	ui.BlockPadding = settings.BlockPadding
	app.Route("/", components.NewHomePage())
	app.Route("/inbound_api", components.NewInboundAPIPage())
	app.Route("/business_api", components.NewBusinessAPIPage())
	app.Route("/monitor", components.NewMonitorPage())
}

func Run() {
	app.RunWhenOnBrowser()
}
