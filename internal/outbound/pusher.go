package outbound

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/outbound/mqtt"
	"github.com/L-LYR/pns/internal/service/cache"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/util"
)

type Pusher interface {
	Handle(context.Context, *model.PushTask) error
}

var _ Pusher = (*mqtt.Client)(nil)

type _PusherManager struct {
	MQTTPushers map[int]Pusher
}

var (
	PusherManager = &_PusherManager{
		MQTTPushers: make(map[int]Pusher),
	}
)

func (p *_PusherManager) MustRegisterPushers(ctx context.Context) {
	brokerConfig := config.MustLoadMQTTBrokerConfig(ctx, "mqtt")
	cache.Config.RangePusherConfig(
		model.MQTTPusher,
		func(appId int, pusherConfig model.PusherConfig) error {
			p.MQTTPushers[appId] = mqtt.MustNewPusher(
				ctx,
				appId,
				pusherConfig.(*model.MQTTConfig),
				brokerConfig,
			)
			util.GLog.Infof(ctx, "Success to initialize mqtt pusher of app %d", appId)
			return nil
		},
	)
	util.GLog.Infof(ctx, "Success to initialize %d mqtt pusher", len(p.MQTTPushers))
}

func (p *_PusherManager) Handle(ctx context.Context, task *model.PushTask, pusher model.PusherType) error {
	log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "task handle", "success")
	var err error
	switch pusher {
	case model.MQTTPusher:
		pusher, ok := p.MQTTPushers[task.Target.App.ID]
		if !ok {
			err = errors.New("pusher not found")
			break
		}
		err = pusher.Handle(ctx, task)
	default:
		panic("unreachable")
	}
	if err != nil {
		log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "push", fmt.Sprintf("Error: %s, fail", err.Error()))
		return err
	}
	log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "push", "success")
	return nil
}
