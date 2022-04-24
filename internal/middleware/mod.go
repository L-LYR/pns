package middleware

import "github.com/gogf/gf/v2/net/ghttp"

var (
	GeneralMiddlewares = []ghttp.HandlerFunc{
		DebugHandler,
		LoggingHandler,
		CommonResponseHandler,
	}
)
