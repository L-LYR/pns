package admin

import (
	"context"

	"github.com/L-LYR/pns/internal/admin/frontend"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

/*
	Only this file can depend on goframe.
*/

const (
	_ServerName = "admin"
)

func MustRegisterRouters(ctx context.Context) *ghttp.Server {
	frontend.MustRegisterFrontendRouters()

	s := g.Server(_ServerName)
	if err := s.SetConfig(*config.AdminServerConfig()); err != nil {
		util.GLog.Panicf(ctx, "Fail to initialize admin server, because %s", err.Error())
	}
	s.BindHandler("/**", func(r *ghttp.Request) {
		frontend.Admin.ServeHTTP(r.Response.Writer, r.Request)
	})

	return s
}
