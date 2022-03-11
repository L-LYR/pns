package main

import (
	"github.com/L-LYR/pns/internal/module"
	"github.com/L-LYR/pns/internal/module/inbound/controller"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	initConfig()
	s := g.Server()
	registerRouter(s)
	module.MustInit()
	s.Run()

	module.MustClean()
}

func initConfig() {
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("config.toml")
}

func registerRouter(s *ghttp.Server) {
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/target", controller.UpsertTarget)
		group.PATCH("/target", controller.UpsertTarget)
		group.PUT("/target", controller.UpsertTarget)
	})
	s.BindHandler("/metrics", func(r *ghttp.Request) {
		promhttp.Handler().ServeHTTP(r.Response.Writer, r.Request)
	})
}
