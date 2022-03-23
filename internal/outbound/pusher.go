package outbound

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/outbound/mqtt"
)

type _Pusher interface {
	Handle(context.Context, *model.PushTask) error
}

var _ _Pusher = (*mqtt.Pusher)(nil)

func _MustNewPusher(
	ctx context.Context,
	appId int,
	t model.PusherType,
) _Pusher {
	switch t {
	case model.MQTTPusher:
		return mqtt.MustNewDefaultPusher(ctx, appId)
	default:
		panic("unreachable")
	}
}

type _PusherManager struct {
	MQTTPushers map[int]_Pusher
}

var (
	_Manager = &_PusherManager{
		MQTTPushers: make(map[int]_Pusher),
	}
)

func MustRegisterPushers(ctx context.Context) {
	appId := 12345
	_Manager.MQTTPushers[appId] = _MustNewPusher(ctx, appId, model.MQTTPusher)
}

func Handle(ctx context.Context, task *model.PushTask, pusher model.PusherType) error {
	switch pusher {
	case model.MQTTPusher:
		pusher, ok := _Manager.MQTTPushers[task.Target.App.ID]
		if !ok {
			return errors.New("pusher not found")
		}
		return pusher.Handle(ctx, task)
	default:
		panic("unreachable")
	}
}
