package main

import (
	"github.com/L-LYR/pns/internal/module/event_queue"
	"github.com/L-LYR/pns/internal/module/inbound"
	"github.com/L-LYR/pns/internal/module/monitor"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

const (
	ConfigFileName = "config.toml"
)

func main() {
	initConfig()

	event_queue.MustInit()
	monitor.MustRegisterMetrics()
	inbound.MustRegisterRouters().Run()

	event_queue.MustShutdown()
}

func initConfig() {
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(ConfigFileName)
}
