package outbound

import (
	"context"
	"fmt"
	"sync"
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

func _MustNewPusher(
	ctx context.Context,
	appId int,
	t model.PusherType,
	pusherConfig model.PusherConfig,
) Pusher {
	switch t {
	case model.MQTTPusher:
		return mqtt.MustNewPusher(
			ctx,
			appId,
			pusherConfig.(*model.MQTTConfig),
			config.MustLoadMQTTBrokerConfig(ctx),
		)
	default:
		panic("unreachable")
	}
}

type _PusherManager struct {
	pusherMutex sync.RWMutex
	pusherType  model.PusherType
	pushers     map[int]Pusher
}

var (
	MQTTPusherManager = &_PusherManager{
		pusherType: model.MQTTPusher,
		pushers:    make(map[int]Pusher),
	}
)

func MustInitialize(ctx context.Context) {
	MQTTPusherManager.MustRegisterPushers(ctx, model.MQTTPusher)
}

// This function is used in initialization stage which is not concurrent-safe.
func (p *_PusherManager) MustRegisterPushers(ctx context.Context, pusherType model.PusherType) {
	if p.pusherType != pusherType {
		panic("unmatched pusher type")
	}
	cache.Config.RangePusherConfig(
		pusherType,
		func(appId int, config model.PusherConfig) error {
			p.pushers[appId] = _MustNewPusher(ctx, appId, p.pusherType, config)
			util.GLog.Infof(ctx, "Success to initialize %s pusher for app %d", pusherType.Name(), appId)
			return nil
		},
	)
	util.GLog.Infof(ctx, "%d %s pushers are running", len(p.pushers), pusherType.Name())
}

// This function will try to new a pusher.
func (p *_PusherManager) _GetPusher(ctx context.Context, appId int, pusherType model.PusherType) (Pusher, bool) {
	p.pusherMutex.RLock()
	if pusher, ok := p.pushers[appId]; ok {
		p.pusherMutex.RUnlock()
		return pusher, true
	}
	p.pusherMutex.RUnlock()
	if pusher, ok := p._TryAddPusher(ctx, appId, pusherType); ok {
		util.GLog.Infof(ctx, "Success to initialize %s pusher for app %d", pusherType.Name(), appId)
		return pusher, true
	}
	util.GLog.Infof(ctx, "Fail to initialize %s pusher for app %d", pusherType.Name(), appId)
	return nil, false
}

func (p *_PusherManager) _TryAddPusher(ctx context.Context, appId int, pusherType model.PusherType) (Pusher, bool) {
	config, ok := cache.Config.GetPusherConfigByAppId(appId, pusherType)
	util.GLog.Infof(ctx, "Try to add pusher")
	if !ok {
		return nil, false
	}
	p.pusherMutex.Lock()
	defer p.pusherMutex.Unlock()
	if pusher, ok := p.pushers[appId]; ok {
		return pusher, true
	}
	pusher := _MustNewPusher(ctx, appId, pusherType, config)
	p.pushers[appId] = pusher
	return pusher, true
}

// Try to re-put task into queue, but if it is failed, task status will be set as dead
func (p *_PusherManager) _ReputTask(ctx context.Context, task *model.PushTask, pusherType model.PusherType) {
	task.Retry++
	if task.Retry > 3 {
		log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "retry", "dead")
		return
	}
	if err := PutPushTaskEvent(ctx, task, pusherType); err != nil {
		log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "retry", fmt.Sprintf("Error: %s, fail", err.Error()))
		return
	}
	log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "retry", "success")
}

func (p *_PusherManager) Handle(ctx context.Context, task *model.PushTask, pusherType model.PusherType) error {
	pusher, ok := p._GetPusher(ctx, task.Target.App.ID, pusherType)
	if !ok {
		util.GLog.Warningf(ctx, "No mqtt pusher for app %d", task.Target.App.ID)
		p._ReputTask(ctx, task, pusherType)
		return nil
	} else if err := pusher.Handle(ctx, task); err != nil {
		log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "push", fmt.Sprintf("Error: %s, fail", err.Error()))
		return err
	}
	log.PutLogEvent(ctx, task.LogMeta(), time.Now().UnixMilli(), "push", "success")
	return nil
}
