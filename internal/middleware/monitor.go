package middleware

import (
	"strconv"
	"time"

	"github.com/L-LYR/pns/internal/monitor"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

func LoggingHandler(r *ghttp.Request) {
	start := time.Now()

	uri := r.RequestURI

	method := r.Method

	r.Middleware.Next()

	duration := time.Since(start)

	if err := r.GetError(); err != nil {
		errCode := strconv.FormatInt(int64(gerror.Code(err).Code()), 10)
		monitor.RequestGenericDuration.
			WithLabelValues(uri, "failure", errCode).Observe(duration.Seconds())
		monitor.RequestGenericCounter.
			WithLabelValues(uri, "failure", errCode).Inc()
	} else {
		monitor.RequestGenericDuration.
			WithLabelValues(uri, "success", "0").Observe(duration.Seconds())
		monitor.RequestGenericCounter.
			WithLabelValues(uri, "success", "0").Inc()
	}

	util.GLog.Infof(r.GetCtx(), "%s %s duration: %s", method, uri, duration.String())
}
