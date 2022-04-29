package main

import (
	"context"

	"github.com/L-LYR/pns/internal/admin"
	"github.com/L-LYR/pns/internal/bizapi"
	"github.com/L-LYR/pns/internal/bizcore"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/inbound"
	"github.com/L-LYR/pns/internal/monitor"
	"github.com/L-LYR/pns/internal/outbound"
	"github.com/L-LYR/pns/internal/service"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/util"
	"github.com/L-LYR/pns/internal/validator"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	/* context & config */
	ctx := GetStartContext()
	config.MustLoad(ctx)
	/* individual modules */
	validator.MustRegisterRules(ctx)
	monitor.MustRegisterMetrics(ctx)
	bizcore.MustInitialize(ctx)
	service.MustInitialize(ctx)
	outbound.MustInitialize(ctx)
	/* event queue */
	EventQueueRegister()
	event_queue.EventQueueManager.MustStart(ctx)
	/* servers */
	inbound.MustRegisterRouters(ctx).Start()
	bizapi.MustRegisterRouters(ctx).Start()
	admin.MustRegisterRouters(ctx).Start()
	/* debug pprof server */
	StartPProf()
	/* block */
	g.Wait()
	/* clean up */
	service.MustShutdown(ctx)
	outbound.MustShutdown(ctx)
	event_queue.EventQueueManager.MustShutdown(ctx)
}

func GetStartContext() context.Context {
	return context.Background()
}

func StartPProf() {
	if util.DebugOn() {
		go ghttp.StartPProfServer(10085)
	}
}

func EventQueueRegister() {
	// temporarily, add a queue between bizapi and bizcore
	event_queue.EventQueueManager.MustRegister(
		config.TaskValidationEventConsumerConfig(),
		bizcore.TaskValidationEventConsumer,
	)
	event_queue.EventQueueManager.MustRegister(
		config.DirectPushTaskEventConsumerConfig(),
		outbound.PushTaskEventConsumer,
	)
	event_queue.EventQueueManager.MustRegister(
		config.RangePushTaskEventConsumerConfig(),
		outbound.PushTaskEventConsumer,
	)
	event_queue.EventQueueManager.MustRegister(
		config.BroadcastPushTaskEventConsumerConfig(),
		outbound.PushTaskEventConsumer,
	)
	event_queue.EventQueueManager.MustRegister(
		config.PushLogEventConsumerConfig(),
		log.PushLogEventConsumer,
	)
}
