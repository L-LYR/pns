package outbound

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/outbound/mqtt"
)

type Pusher interface {
	Handle(context.Context, *model.PushTask) error
}

var _ Pusher = (*mqtt.Client)(nil)

func MustNewPusher(
	ctx context.Context,
	appId int,
	t model.PusherType,
) Pusher {
	switch t {
	case model.MQTTPusher:
		return mqtt.MustNewPusher(ctx, appId, config.MustLoadMQTTBrokerConfig(ctx, "mqtt"))
	default:
		panic("unreachable")
	}
}

type _PusherManager struct {
	MQTTPushers map[int]Pusher
}

var (
	PusherManager = &_PusherManager{
		MQTTPushers: make(map[int]Pusher),
	}
)

func (p *_PusherManager) MustRegisterPushers(ctx context.Context) {
	appId := 12345
	p.MQTTPushers[appId] = MustNewPusher(ctx, appId, model.MQTTPusher)
}

func (p *_PusherManager) Handle(ctx context.Context, task *model.PushTask, pusher model.PusherType) error {
	switch pusher {
	case model.MQTTPusher:
		pusher, ok := p.MQTTPushers[task.Target.App.ID]
		if !ok {
			return errors.New("pusher not found")
		}
		return pusher.Handle(ctx, task)
	default:
		panic("unreachable")
	}
}
