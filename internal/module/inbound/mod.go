package inbound

import (
	"context"

	"github.com/L-LYR/pns/internal/module/inbound/controller"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ServerName       = "inbound"
	ServerConfigName = "server.inbound"
)

func MustRegisterRouters() *ghttp.Server {
	s := g.Server(ServerName)
	s.SetConfigWithMap(g.Cfg().MustGet(context.Background(), ServerConfigName).Map())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/target", controller.UpsertTarget)
		group.PATCH("/target", controller.UpsertTarget)
		group.PUT("/target", controller.UpsertTarget)
	})
	s.BindHandler("/metrics", func(r *ghttp.Request) {
		promhttp.Handler().ServeHTTP(r.Response.Writer, r.Request)
	})
	return s
}
