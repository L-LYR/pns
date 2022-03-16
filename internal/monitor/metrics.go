package monitor

import "github.com/prometheus/client_golang/prometheus"

// TODO: add a metric manager
var (
	UpsertTargetEventMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upsert_target",
		},
		[]string{"type", "result"},
	)
)

func MustRegisterMetrics() {
	prometheus.MustRegister(UpsertTargetEventMetric)
}
