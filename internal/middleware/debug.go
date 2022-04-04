package middleware

import (
	"time"

	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/net/ghttp"
)

// This middelware must be put at the very beginning.
func DebugHandler(r *ghttp.Request) {
	// add whatever you want
	r.Middleware.Next()
}

func LoggingHandler(r *ghttp.Request) {
	start := time.Now()

	uri := r.RequestURI
	method := r.Method

	r.Middleware.Next()

	duration := time.Since(start)

	util.GLog.Infof(r.GetCtx(), "%s %s duration: %s", method, uri, duration.String())
}
