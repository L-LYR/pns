package main

import (
	"context"

	"github.com/L-LYR/pns/internal/admin"
	"github.com/L-LYR/pns/internal/bizapi"
	"github.com/L-LYR/pns/internal/constdef"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/inbound"
	"github.com/L-LYR/pns/internal/monitor"
	"github.com/L-LYR/pns/internal/outbound"
	"github.com/L-LYR/pns/internal/service"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

func main() {
	/* context & config */
	ctx := GetStartContext()
	LoadGlobalConfig()

	/* individual modules */
	monitor.MustRegisterMetrics()
	service.MustInitialize(ctx)
	outbound.PusherManager.MustRegisterPushers(ctx)

	/* event queue */
	event_queue.EventQueueManager.MustRegister(
		constdef.PushEventTopic,
		g.Cfg().MustGet(ctx, "module.event_queue.push_event_consumer").Uint(),
		outbound.PushEventConsumer,
	)
	event_queue.EventQueueManager.MustRegister(
		constdef.TargetEventTopic,
		g.Cfg().MustGet(ctx, "module.event_queue.target_event_consumer").Uint(),
		target.TargetEventConsumer,
	)
	event_queue.EventQueueManager.MustStart()
	/* servers */
	inbound.MustRegisterRouters(ctx).Start()
	bizapi.MustRegisterRouters(ctx).Start()
	admin.MustRegisterRouters(ctx).Start()

	g.Wait()

	/* clean up */
	event_queue.EventQueueManager.MustShutdown()
	service.MustShutdown(ctx)
}

func GetStartContext() context.Context {
	return context.Background()
}

func LoadGlobalConfig() {
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("config.toml")
}
