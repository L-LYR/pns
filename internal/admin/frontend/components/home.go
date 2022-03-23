package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Home struct {
	app.Compo

	name string
}

func NewHomePage() *Home {
	return &Home{}
}

func (h *Home) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello, "),
			app.If(h.name != "", app.Text(h.name)).
				Else(app.Text("World!")),
		),
		app.P().Body(
			app.Input().
				Type("text").
				Value(h.name).
				Placeholder("What is your name?").
				AutoFocus(true).
				OnChange(h.ValueTo(&h.name)),
		),
	)
}
