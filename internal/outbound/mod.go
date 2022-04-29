package outbound

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

var (
	MQTTPusherManager = &_PusherManager{
		pusherType: model.MQTTPusher,
		pushers:    make(map[int]Pusher),
	}
)

func MustInitialize(ctx context.Context) {
	MQTTPusherManager.MustRegisterPushers(ctx, model.MQTTPusher)
}

func MustShutdown(ctx context.Context) {
	MQTTPusherManager.MustShutdown(ctx)
}
