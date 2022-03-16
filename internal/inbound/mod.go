package inbound

import (
	"context"

	"github.com/L-LYR/pns/internal/inbound/controller"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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
		group.GET("/target", controller.QueryTarget)
		group.POST("/target", controller.CreateTarget)
		group.PATCH("/target", controller.UpdateTarget)
	})
	s.BindHandler("/metrics", func(r *ghttp.Request) {
		promhttp.Handler().ServeHTTP(r.Response.Writer, r.Request)
	})
	return s
}
