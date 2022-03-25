package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

type basePage struct {
	app.Compo

	class   string
	content []app.UI
	index   []app.UI

	updateAvailable bool
}

func newBasePage() *basePage {
	return &basePage{}
}

func (p *basePage) Class(v string) *basePage {
	p.class = app.AppendClass(p.class, v)
	return p
}

func (p *basePage) Index(v ...app.UI) *basePage {
	p.index = app.FilterUIElems(v...)
	return p
}

func (p *basePage) Content(v ...app.UI) *basePage {
	p.content = app.FilterUIElems(v...)
	return p
}

func (p *basePage) OnNav(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
	ctx.Defer(func(ctx app.Context) {
		id := ctx.Page().URL().Fragment
		if id == "" {
			id = "page-top"
		}
		ctx.ScrollTo(id)
	})
}

func (p *basePage) OnAppUpdate(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

func (p *basePage) Render() app.UI {
	return ui.Shell().
		Class("fill").
		Class("background").
		HamburgerMenu(
			newMenu().Class("fill").Class("menu-hamburger-background"),
		).
		Menu(
			newMenu().Class("fill"),
		).
		Content(
			app.Main().Class("fill").Body(
				app.Range(p.content).Slice(func(i int) app.UI {
					return p.content[i]
				}),
			),
		)
}
