package components

import (
	"strings"

	"github.com/L-LYR/pns/internal/admin/frontend/settings"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

type menu struct {
	app.Compo

	class string
}

func newMenu() *menu {
	return &menu{}
}

func (m *menu) Class(v string) *menu {
	m.class = app.AppendClass(m.class, v)
	return m
}

func (m *menu) Render() app.UI {
	linkClass := strings.Join([]string{"link", "heading", "fit", "unselectable"}, " ")
	focus := func(path string) string {
		if app.Window().URL().Path == path {
			return "focus"
		}
		return ""
	}
	return ui.Scroll().
		Class("menu").
		Class(m.class).
		HeaderHeight(settings.HeaderHeight).
		Header(
			ui.Stack().Class("fill").Middle().Content(
				app.Header().Body(
					app.A().Class("heading").Class("app-title").Href("/").Text("PNS Admin"),
				),
			),
		).
		Content(
			app.Nav().Body(
				app.Div().Class("separator"),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Inbound API").Href("/inbound_api").Class(focus("/inbound_api")),
				ui.Link().Class(linkClass).Label("Business API").Href("/business_api").Class(focus("/business_api")),
				ui.Link().Class(linkClass).Label("Monitor").Href("/monitor").Class(focus("/monitor")),
				app.Div().Class("separator"),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				app.Div().Class("separator"),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				app.Div().Class("separator"),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				app.Div().Class("separator"),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				ui.Link().Class(linkClass).Label("Home").Href("/").Class(focus("/")),
				app.Div().Class("separator"),
			),
		)
}
