package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type HomePage struct {
	app.Compo
}

func NewHomePage() *HomePage {
	return &HomePage{}
}

func (h *HomePage) Render() app.UI {
	return newBasePage()
}
