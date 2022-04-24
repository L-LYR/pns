package event_queue

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

type Event interface {
	GetCtx() context.Context
}

var (
	_ Event = (*model.PushTaskEvent)(nil)
	_ Event = (*model.PushLogEvent)(nil)
)

type EventQueue interface {
	Start(context.Context)
	Put(string, Event) error
	Subscribe(string) (<-chan Event, error)
	Shutdown(context.Context)
}

var (
	_ EventQueue = (*_InMemoryEventQueue)(nil)
)

func _MustNewEventQueue(topics ...string) EventQueue {
	return _MustNewInMemoryEventQueue(topics)
}
