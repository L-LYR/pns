package event_queue

import "context"

// Event must contain the context of the request.
type Event interface {
	GetCtx() context.Context
}
type EventQueue interface {
	Start()
	Put(string, Event) error
	Subscribe(string) (<-chan Event, error)
	Shutdown()
}

func NewEventQueue(topics []string) EventQueue {
	return newInMemoryEventQueue(topics)
}
