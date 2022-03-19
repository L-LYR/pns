package inbound

import (
	"context"

	"github.com/L-LYR/pns/internal/inbound/controller"
	"github.com/L-LYR/pns/internal/middleware"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/protocol/goai"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ServerName       = "inbound"
	ServerConfigName = "server.inbound"
)

func MustRegisterRouters(ctx context.Context) *ghttp.Server {
	s := g.Server(ServerName)
	s.SetConfigWithMap(g.Cfg().MustGet(ctx, ServerConfigName).Map())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			middleware.DebugHandler,
			middleware.CommonResponseHandler,
		)
		group.Bind(controller.Target)
	})
	s.BindHandler("/metrics", func(r *ghttp.Request) {
		promhttp.Handler().ServeHTTP(r.Response.Writer, r.Request)
	})
	s.SetOpenApiPath("/api")
	s.SetSwaggerPath("/swagger")
	openapi := s.GetOpenApi()
	openapi.Config.CommonResponse = middleware.CommonResponse{}
	openapi.Config.CommonResponseDataField = `Payload`
	openapi.Info.Title = `Push Inbound API Reference`
	openapi.Info.Description = `This is a description.`
	openapi.Tags = &goai.Tags{
		{
			Name:        "Inbound",
			Description: "This is a description.",
		},
	}
	return s
}
