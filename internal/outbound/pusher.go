package outbound

import (
	"context"
	"errors"
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
	Handle(context.Context, model.PushTask) error
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
			config.MQTTBrokerConfig(),
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
func (p *_PusherManager) _GetPusher(ctx context.Context, appId int) (Pusher, bool) {
	p.pusherMutex.RLock()
	if pusher, ok := p.pushers[appId]; ok {
		p.pusherMutex.RUnlock()
		return pusher, true
	}
	p.pusherMutex.RUnlock()
	if pusher, ok := p._TryAddPusher(ctx, appId); ok {
		util.GLog.Infof(ctx, "Success to initialize %s pusher for app %d", p.pusherType.Name(), appId)
		return pusher, true
	}
	util.GLog.Infof(ctx, "Fail to initialize %s pusher for app %d", p.pusherType.Name(), appId)
	return nil, false
}

func (p *_PusherManager) _TryAddPusher(ctx context.Context, appId int) (Pusher, bool) {
	config, ok := cache.Config.GetPusherConfigByAppId(appId, p.pusherType)
	util.GLog.Infof(ctx, "Try to add pusher")
	if !ok {
		return nil, false
	}
	p.pusherMutex.Lock()
	defer p.pusherMutex.Unlock()
	if pusher, ok := p.pushers[appId]; ok {
		return pusher, true
	}
	pusher := _MustNewPusher(ctx, appId, p.pusherType, config)
	p.pushers[appId] = pusher
	return pusher, true
}

// Try to re-put task into queue, but if it is failed, task status will be set as dead
func (p *_PusherManager) _ReputTask(ctx context.Context, task model.PushTask, pusherType model.PusherType) error {
	meta := task.GetLogMeta()
	task.GetMeta().SetRetry()
	if task.CanRetry() {
		if err := log.PutTaskLog(ctx, meta, "retry", "failure"); err != nil {
			util.GLog.Warningf(ctx, "Fail to set task log, err = %+v", err)
		}
		return errors.New("fail to retry")
	}
	if err := PutPushTaskEvent(ctx, task); err != nil {
		util.GLog.Warningf(ctx, "Task %d fail to reput task in event queue, retry", task.GetID())
		return p._ReputTask(ctx, task, pusherType)
	}
	if err := log.PutTaskLog(ctx, meta, "retry", "success"); err != nil {
		util.GLog.Warningf(ctx, "Fail to set task log, err = %+v", err)
	}
	return nil
}

func (p *_PusherManager) Handle(ctx context.Context, task model.PushTask) error {
	taskMeta := task.GetMeta()
	if taskMeta.OnHandle() {
		taskMeta.SetHandleTime(time.Now())
	}
	logMeta := task.GetLogMeta()
	pusher, ok := p._GetPusher(ctx, task.GetAppId())
	if !ok {
		util.GLog.Warningf(ctx, "No %s pusher for app %d", task.GetPusher().Name(), task.GetAppId())
		return p._ReputTask(ctx, task, p.pusherType)
	}
	if err := log.PutTaskEntry(ctx, logMeta); err != nil {
		util.GLog.Warningf(ctx, "Fail to add task list entry, meta:%+v", logMeta)
	}
	if err := pusher.Handle(ctx, task); err != nil {
		if err := p._ReputTask(ctx, task, p.pusherType); err != nil {
			taskMeta.SetFailure()
			return err
		} else { // retrying
			return nil
		}
	}
	taskMeta.SetSuccess()
	return nil
}
