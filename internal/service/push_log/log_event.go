package log

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
)

func LogEventConsumer(e event_queue.Event) error {
	le, ok := e.(*model.PushLogEvent)
	if !ok {
		return errors.New("not LogEvent")
	}
	return PutTaskLog(le.GetCtx(), le.GetEntry())
}

func PutLogEvent(
	ctx context.Context,
	meta *model.PushLogMeta,
	timestamp int64,
	where string,
	hint string,
) {
	if err := event_queue.EventQueueManager.Put(
		config.LogEventTopic(),
		&model.PushLogEvent{
			Ctx: ctx,
			Entry: &model.LogEntry{
				LogBase: &model.LogBase{
					Meta:  meta,
					T:     timestamp,
					Where: where,
				},
				Hint: hint,
			},
		},
	); err != nil {
		util.GLog.Errorf(
			ctx,
			"Fail to put log event, meta: %+v, timestamp: %d, where: %s, hint: %s",
			meta, timestamp, where, hint,
		)
	}
}
