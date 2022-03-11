package module

import (
	"github.com/L-LYR/pns/internal/module/event_manager"
	"github.com/L-LYR/pns/internal/module/monitor"
)

func MustInit() {
	monitor.RegisterMetrics()
	event_manager.MustStart()
}

func MustClean() {
	event_manager.MustShutdown()
}
