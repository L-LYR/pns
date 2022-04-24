package main

import (
	"context"

	"github.com/L-LYR/pns/internal/admin"
	"github.com/L-LYR/pns/internal/bizapi"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/inbound"
	"github.com/L-LYR/pns/internal/monitor"
	"github.com/L-LYR/pns/internal/outbound"
	"github.com/L-LYR/pns/internal/service"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/validator"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	/* context & config */
	ctx := GetStartContext()
	config.MustLoadTaskDefaultConfig(ctx)
	/* individual modules */
	validator.MustRegisterRules(ctx)
	monitor.MustRegisterMetrics(ctx)
	service.MustInitialize(ctx)
	outbound.MustInitialize(ctx)
	/* event queue */
	event_queue.EventQueueManager.MustRegister(
		config.MustLoadConsumerConfig(ctx, "push_event_consumer"),
		outbound.PushTaskEventConsumer,
	)
	event_queue.EventQueueManager.MustRegister(
		config.MustLoadConsumerConfig(ctx, "log_event_consumer"),
		log.PushLogEventConsumer,
	)
	event_queue.EventQueueManager.MustStart(ctx)
	/* servers */
	inbound.MustRegisterRouters(ctx).Start()
	bizapi.MustRegisterRouters(ctx).Start()
	admin.MustRegisterRouters(ctx).Start()
	g.Wait()
	/* clean up */
	event_queue.EventQueueManager.MustShutdown(ctx)
	service.MustShutdown(ctx)
}

func GetStartContext() context.Context {
	return context.Background()
}
