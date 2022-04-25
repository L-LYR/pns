package outbound

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/monitor"
	log "github.com/L-LYR/pns/internal/service/push_log"
	"github.com/L-LYR/pns/internal/util"
)

func PushTaskEventConsumer(e event_queue.Event) error {
	pe, ok := e.(*model.PushTaskEvent)
	if !ok {
		return errors.New("not PushEvent")
	}

	ctx, task := pe.GetCtx(), pe.GetTask()
	logMeta := task.GetLogMeta()

	if err := log.PutTaskLog(ctx, logMeta, "task handle", "success"); err != nil {
		util.GLog.Warningf(ctx, "Fail to set task log, err = %+v", err)
	}

	task.GetMeta().SetOnHandle()

	var err error
	taskHint := "success"
	switch pe.GetTask().GetPusher() {
	case model.MQTTPusher:
		err = MQTTPusherManager.Handle(ctx, task)
	default:
		panic("unreachable")
	}

	if !task.GetMeta().IsDone() {
		return nil
	}
	task.GetMeta().SetEndTime(time.Now())

	if err != nil {
		taskHint = "failure"
	}

	if err := log.PutTaskLog(ctx, logMeta, "task done", taskHint); err != nil {
		util.GLog.Warningf(ctx, "Fail to set task log, err = %+v", err)
	}

	taskTypeName := task.GetType().Name()
	monitor.PushTaskCounter.
		WithLabelValues(taskTypeName, "done", taskHint).Inc()
	monitor.PushTaskDuration.
		WithLabelValues(taskTypeName, "total", taskHint).Observe(
		task.GetMeta().TotalDuration().Seconds(),
	)
	monitor.PushTaskDuration.
		WithLabelValues(taskTypeName, "validation", taskHint).Observe(
		task.GetMeta().ValidationDuration().Seconds(),
	)
	monitor.PushTaskDuration.
		WithLabelValues(taskTypeName, "handle", taskHint).Observe(
		task.GetMeta().HandleDuration().Seconds(),
	)

	util.GLog.Infof(ctx, "Task %d %s", task.GetID(), taskHint)

	return err
}

func PutPushTaskEvent(ctx context.Context, task model.PushTask) error {
	return event_queue.EventQueueManager.Put(
		config.PushTaskEventTopic(),
		&model.PushTaskEvent{
			Ctx:  ctx,
			Task: task,
		},
	)
}
