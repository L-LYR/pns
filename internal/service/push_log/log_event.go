package log

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
)

func LogEventConsumer(e event_queue.Event) error {
	le, ok := e.(*model.LogEvent)
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
) error {
	return event_queue.EventQueueManager.Put(
		config.LogEventTopic(),
		&model.LogEvent{
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
	)
}
