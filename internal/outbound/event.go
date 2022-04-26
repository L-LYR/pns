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
		return errors.New("not PushTaskEvent")
	}
	ctx, task := pe.GetCtx(), pe.GetTask()
	logMeta, taskMeta, taskTypeName :=
		task.GetLogMeta(), task.GetMeta(), task.GetType().Name()

	monitor.PushTaskDuration.
		WithLabelValues(taskTypeName, "task pending since validation", "-").Observe(
		time.Since(taskMeta.GetValidationTime()).Seconds(),
	)

	log.PutTaskLogEvent(ctx, logMeta, "task handle", "success")

	taskMeta.SetOnHandle()

	var err error
	taskHint := "success"
	switch pe.GetTask().GetPusher() {
	case model.MQTTPusher:
		err = MQTTPusherManager.Handle(ctx, task)
	default:
		panic("unreachable")
	}

	if !taskMeta.IsDone() {
		return nil
	}
	taskMeta.SetEndTime(time.Now())

	if err != nil {
		taskHint = "failure"
	}

	log.PutTaskLogEvent(ctx, logMeta, "task done", taskHint)

	monitor.PushTaskCounter.
		WithLabelValues(taskTypeName, "outbound", taskHint).Inc()

	monitor.PushTaskDuration.
		WithLabelValues(taskTypeName, "handle", taskHint).Observe(
		taskMeta.GetEndTime().Sub(taskMeta.GetHandleTime()).Seconds(),
	)

	monitor.PushTaskDuration.
		WithLabelValues(taskTypeName, "total", taskHint).Observe(
		taskMeta.GetEndTime().Sub(taskMeta.GetCreationTime()).Seconds(),
	)

	util.GLog.Infof(ctx, "Task %d %s", task.GetID(), taskHint)

	return err
}

func PutPushTaskEvent(ctx context.Context, task model.PushTask) error {
	e := &model.PushTaskEvent{
		Ctx:  ctx,
		Task: task,
	}
	switch task.GetType() {
	case model.BroadcastPush:
		return event_queue.EventQueueManager.Put(config.BroadcastPushTaskEventTopic(), e)
	case model.DirectPush:
		return event_queue.EventQueueManager.Put(config.DirectPushTaskEventTopic(), e)
	default:
		return errors.New("unknown task type")
	}
}
