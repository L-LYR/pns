package bizcore

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/monitor"
	"github.com/L-LYR/pns/internal/outbound"
	log "github.com/L-LYR/pns/internal/service/push_log"
)

func TaskValidationEventConsumer(e event_queue.Event) error {
	pe, ok := e.(*model.PushTaskEvent)
	if !ok {
		return errors.New("not PushTaskEvent")
	}
	ctx, task := pe.GetCtx(), pe.GetTask()
	logMeta, taskMeta := task.GetLogMeta(), task.GetMeta()

	monitor.PushTaskDuration.
		WithLabelValues(task.GetType().Name(), "task pending since creation", "-").Observe(
		time.Since(taskMeta.GetCreationTime()).Seconds(),
	)

	taskMeta.SetPending()
	taskMeta.SetValidationTime(time.Now())

	taskHint := "success"
	err := _Validate(ctx, task)
	if err != nil {
		taskHint = "failure"
	}

	log.PutTaskLogEvent(ctx, logMeta, "validation", taskHint)

	monitor.PushTaskCounter.
		WithLabelValues(task.GetType().Name(), "validation", taskHint).Inc()

	monitor.PushTaskDuration.
		WithLabelValues(task.GetType().Name(), "validation", taskHint).Observe(
		time.Since(taskMeta.GetValidationTime()).Seconds(),
	)

	return err
}

func _Validate(ctx context.Context, task model.PushTask) error {
	if err := Execute(
		map[string]interface{}{
			"ctx":  ctx,
			"task": task,
		},
	); err != nil {
		return err
	}
	if err := outbound.PutPushTaskEvent(ctx, task); err != nil {
		return err
	}
	return nil
}

func PutTaskValidationEvent(ctx context.Context, task model.PushTask) error {
	return event_queue.EventQueueManager.Put(
		config.TaskValidationEventTopic(),
		&model.PushTaskEvent{
			Ctx:  ctx,
			Task: task,
		},
	)
}
