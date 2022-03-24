package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type embedPage struct {
	app.Compo

	name string
	url  string
}

func newEmbedPage() *embedPage {
	return &embedPage{}
}

func (p *embedPage) Src(url string) *embedPage {
	p.url = url
	return p
}

func (p *embedPage) ID(v string) *embedPage {
	p.name = v
	return p
}

func (p *embedPage) Render() app.UI {
	return newBasePage().Content(
		app.IFrame().Class("fill").ID(p.name).Src(p.url),
	)
}

type InboundAPIPage struct {
	embedPage
}

func NewInboundAPIPage() *InboundAPIPage {
	return &InboundAPIPage{}
}

func (p *InboundAPIPage) Render() app.UI {
	return newEmbedPage().ID("inbound").Src("http://localhost:10086/swagger")
}

type BusinessAPIPage struct {
	embedPage
}

func NewBusinessAPIPage() *BusinessAPIPage {
	return &BusinessAPIPage{}
}

func (p *BusinessAPIPage) Render() app.UI {
	return newEmbedPage().ID("business").Src("http://localhost:10087/swagger")
}

type MonitorPage struct {
	embedPage
}

func NewMonitorPage() *MonitorPage {
	return &MonitorPage{}
}

func (p *MonitorPage) Render() app.UI {
	return newEmbedPage().ID("monitor").Src("http://localhost:3000")
}
