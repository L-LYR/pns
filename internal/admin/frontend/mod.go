package frontend

import (
	"github.com/L-LYR/pns/internal/admin/frontend/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var (
	Admin = &app.Handler{
		Author:       "HammerLi",
		Name:         "PNS Admin",
		Title:        "PNS Admin",
		LoadingLabel: "PNS Admin",
		Icon: app.Icon{
			Default: "https://avatars.githubusercontent.com/u/45999891?v=4",
		},
		Styles: []string{
			`https://unpkg.com/@patternfly/patternfly@4.96.2/patternfly.css`,
			`https://unpkg.com/@patternfly/patternfly@4.96.2/patternfly-addons.css`,
			`/web/index.css`,
		},
	}
)

func MustRegisterFrontendRouters() {
	app.Route("/", components.NewHomePage())
}

func Run() {
	app.RunWhenOnBrowser()
}
