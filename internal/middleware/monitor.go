package middleware

import (
	"strconv"
	"time"

	"github.com/L-LYR/pns/internal/monitor"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

func LoggingHandler(r *ghttp.Request) {
	start := time.Now()

	req := r.Method + " " + r.RequestURI

	r.Middleware.Next()

	duration := time.Since(start)

	if err := r.GetError(); err != nil {
		errCode := strconv.FormatInt(int64(gerror.Code(err).Code()), 10)
		monitor.RequestGenericDuration.
			WithLabelValues(req, "failure", errCode).Observe(duration.Seconds())
		monitor.RequestGenericCounter.
			WithLabelValues(req, "failure", errCode).Inc()
	} else {
		monitor.RequestGenericDuration.
			WithLabelValues(req, "success", "0").Observe(duration.Seconds())
		monitor.RequestGenericCounter.
			WithLabelValues(req, "success", "0").Inc()
	}
}
