package bizapi

import (
	"context"

	"github.com/L-LYR/pns/internal/bizapi/controller"
	"github.com/L-LYR/pns/internal/middleware"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/protocol/goai"
)

const (
	ServerName       = "bizapi"
	ServerConfigName = "server.bizapi"
)

func MustRegisterRouters(ctx context.Context) *ghttp.Server {
	s := g.Server(ServerName)
	s.SetConfigWithMap(g.Cfg().MustGet(ctx, ServerConfigName).Map())
	// Bind all controller objects
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			middleware.DebugHandler,
			middleware.CommonResponseHandler,
		)
		group.Bind(controller.Push)
	})
	// Register Open API docs
	s.SetOpenApiPath("/api")
	s.SetSwaggerPath("/swagger")
	// TODO: add descriptions
	openapi := s.GetOpenApi()
	openapi.Config.CommonResponse = middleware.CommonResponse{}
	openapi.Config.CommonResponseDataField = `Payload`
	openapi.Info.Title = `Push Business API Reference`
	openapi.Info.Description = `This is a description.`
	openapi.Tags = &goai.Tags{
		{
			Name:        "Push",
			Description: "This is a description.",
		},
	}
	return s
}
