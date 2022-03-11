package monitor

import "github.com/prometheus/client_golang/prometheus"

// TODO: add a metric manager
var (
	UpsertTargetMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upsert_target",
		},
		[]string{"type", "result"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(UpsertTargetMetric)
}
