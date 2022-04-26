package log

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
)

func PushLogEventConsumer(e event_queue.Event) error {
	le, ok := e.(*model.LogEvent)
	if !ok {
		return errors.New("not LogEvent")
	}

	switch le.GetType() {
	case model.PushLog:
		return PutPushLog(le.GetCtx(), le.GetEntry())
	case model.TaskLog:
		return PutTaskLog(le.GetCtx(), le.GetEntry())
	default:
		panic("unreachable")
	}
}

func PutTaskLogEvent(
	ctx context.Context,
	meta *model.LogMeta,
	where string, hint string,
) {
	PutLogEvent(ctx, meta, where, time.Now().UnixMilli(), hint, model.TaskLog)
}

func PutPushLogEvent(
	ctx context.Context,
	meta *model.LogMeta,
	where string,
	when int64,
	hint string,
) {
	PutLogEvent(ctx, meta, where, when, hint, model.PushLog)
}

func PutLogEvent(
	ctx context.Context,
	meta *model.LogMeta,
	where string,
	when int64,
	hint string,
	logEventType model.LogEventType,
) {
	if err := event_queue.EventQueueManager.Put(
		config.PushLogEventTopic(),
		&model.LogEvent{
			Ctx: ctx,
			Entry: &model.LogEntry{
				LogBase: &model.LogBase{
					Meta:  meta,
					T:     when,
					Where: where,
				},
				Hint: hint,
			},
			Type: logEventType,
		},
	); err != nil {
		util.GLog.Errorf(
			ctx,
			"Fail to put log event, meta: %+v, timestamp: %d, where: %s, hint: %s",
			meta, when, where, hint,
		)
	}
}
