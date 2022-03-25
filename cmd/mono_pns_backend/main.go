package main

import (
	"context"

	"github.com/L-LYR/pns/internal/admin"
	"github.com/L-LYR/pns/internal/bizapi"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/inbound"
	"github.com/L-LYR/pns/internal/monitor"
	"github.com/L-LYR/pns/internal/outbound"
	"github.com/L-LYR/pns/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

const (
	ConfigFileName = "config.toml"
)

func main() {
	ctx := context.Background()

	InitGlobalConfig()
	service.MustInit(ctx)
	event_queue.MustInit()
	monitor.MustRegisterMetrics()

	outbound.MustRegisterPushers(ctx)
	inbound.MustRegisterRouters(ctx).Start()
	bizapi.MustRegisterRouters(ctx).Start()
	admin.MustRegisterRouters(ctx).Start()

	g.Wait()

	event_queue.MustShutdown()
	service.MustShutdown(ctx)
}

func InitGlobalConfig() {
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(ConfigFileName)
}
