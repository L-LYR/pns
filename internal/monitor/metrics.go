package monitor

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

// TODO: add a metric manager
var (
	RequestGenericTags = []string{
		"uri",     // uri
		"result",  // success or failure
		"errCode", // errCode
	}
	RequestGenericCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "request"},
		RequestGenericTags,
	)
	RequestGenericDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{Name: "request_duration"},
		RequestGenericTags,
	)
)

var (
	PushTaskTags = []string{
		"type",   // direct broadcast
		"where",  // creation, validation, retry, outbound ......
		"result", // success or failure
	}
	PushTaskCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "push_task"},
		PushTaskTags,
	)
	PushTaskDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{Name: "push_task_duration"},
		PushTaskTags,
	)
)

var (
	EventQueueTags = []string{
		"topic",
	}
	EventQueueLength = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "event_queue_len"},
		EventQueueTags,
	)
)

func MustRegisterMetrics(ctx context.Context) {
	prometheus.MustRegister(
		RequestGenericCounter,
		RequestGenericDuration,
		PushTaskCounter,
		PushTaskDuration,
		EventQueueLength,
	)
}
