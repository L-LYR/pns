package middleware

import (

	"github.com/gogf/gf/v2/net/ghttp"
)

// This middelware must be put at the very beginning.
func DebugHandler(r *ghttp.Request) {
	// add whatever you want
	r.Middleware.Next()
}