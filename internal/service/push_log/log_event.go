package log

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
)

func PushLogEventConsumer(e event_queue.Event) error {
	le, ok := e.(*model.PushLogEvent)
	if !ok {
		return errors.New("not LogEvent")
	}
	return PutPushLog(le.GetCtx(), le.GetEntry())
}

func PutPushLogEvent(
	ctx context.Context,
	hint string,
	base *model.LogBase,
) {
	if err := event_queue.EventQueueManager.Put(
		config.PushLogEventTopic(),
		&model.PushLogEvent{
			Ctx: ctx,
			Entry: &model.LogEntry{
				LogBase: base,
				Hint:    hint,
			},
		},
	); err != nil {
		util.GLog.Errorf(
			ctx,
			"Fail to put log event, meta: %+v, timestamp: %d, where: %s, hint: %s",
			base.Meta, base.T, base.Where, hint,
		)
	}
}
