package task

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/bizcore"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/cache"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
	"golang.org/x/time/rate"
)

type RangePushTaskRunner struct {
	ctx     context.Context
	task    *model.RangePushTask
	limiter *rate.Limiter
}

func NewRangePushTaskRunner(ctx context.Context, task *model.RangePushTask, limiter *rate.Limiter) *RangePushTaskRunner {
	return &RangePushTaskRunner{
		ctx:  ctx,
		task: task,
		// TODO: configurable
		// TODO: range_runner should share limiter with apis
		limiter: limiter,
	}
}

func (r *RangePushTaskRunner) Run() {
	appName, _ := cache.Config.GetAppNameByAppId(r.task.GetAppId())

	cursor, err := target.Cursor(r.ctx, appName, r.task.CursorFilters()...)
	if err != nil {
		util.GLog.Warningf(r.ctx, "Fail to get cursor for range push task %d", r.task.GetID())
		return
	}
	for cursor.Next(r.ctx) {
		target := &model.Target{}
		if err := cursor.Decode(target); err != nil {
			// TODO: remove this
			util.GLog.Warningf(r.ctx, "Fail to decode target, because %s", err.Error())
			continue
		}
		newTask := r.task.Spawn()
		newTask.Target = target
		r.limiter.Wait(r.ctx)
		if err := bizcore.PutTaskValidationEvent(r.ctx, newTask); err != nil {
			util.GLog.Warningf(r.ctx, "Fail to put task, because %s", err.Error())
			if errors.Is(err, event_queue.ErrClose) {
				break
			}
			continue
		}
	}
	util.GLog.Infof(r.ctx, "Range task %d done", r.task.GetID())
	log.PutLogEvent(r.ctx, r.task.GetLogMeta(), model.TaskDone, time.Now().UnixMilli(), "success", model.TaskLog)
}
