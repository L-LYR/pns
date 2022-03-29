package outbound

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/outbound/mqtt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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
		return mqtt.MustNewPusher(ctx, appId, _LoadDefaultMQTTBrokerConfig(ctx))
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

const (
	_PusherConfigName = "module.pusher.mqtt"
)

func _LoadDefaultMQTTBrokerConfig(ctx context.Context) *mqtt.BrokerConfig {
	brokerConfig := &mqtt.BrokerConfig{}
	config := g.Cfg().MustGet(ctx, _PusherConfigName).Map()
	if err := gconv.Struct(config, brokerConfig); err != nil {
		panic("fail to initialize default pusher config")
	}
	return brokerConfig
}
