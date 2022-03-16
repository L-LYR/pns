package target

import "github.com/L-LYR/pns/internal/monitor"

func EmitUpsertTargetEvent(what, result string) {
	monitor.UpsertTargetEventMetric.WithLabelValues(what, result).Inc()
}
