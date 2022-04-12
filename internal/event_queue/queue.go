package event_queue

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

type Event interface {
	GetCtx() context.Context
}

var (
	_ Event = (*model.PushEvent)(nil)
	_ Event = (*model.LogEvent)(nil)
)

type EventQueue interface {
	Start()
	Put(string, Event) error
	Subscribe(string) (<-chan Event, error)
	Shutdown()
}

var (
	_ EventQueue = (*_InMemoryEventQueue)(nil)
)

func _MustNewEventQueue(topics ...string) EventQueue {
	return _MustNewInMemoryEventQueue(topics)
}
