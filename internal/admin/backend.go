package admin

import (
	"context"

	"github.com/L-LYR/pns/internal/admin/frontend"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	_ServerName       = "admin"
	_ServerConfigName = "server.admin"
)

func MustRegisterRouters(ctx context.Context) *ghttp.Server {
	frontend.MustRegisterFrontendRouters()

	s := g.Server(_ServerName)
	s.SetConfigWithMap(g.Cfg().MustGet(ctx, _ServerConfigName).Map())
	s.BindHandler("/**", func(r *ghttp.Request) {
		frontend.Admin.ServeHTTP(r.Response.Writer, r.Request)
	})

	return s
}
