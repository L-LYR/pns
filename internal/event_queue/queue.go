package event_queue

import "context"

// _Event must contain the context of the request.
type _Event interface {
	GetCtx() context.Context
}

var (
	_ _Event = (*_TargetEvent)(nil)
	_ _Event = (*_PushEvent)(nil)
)

type _EventQueue interface {
	Start()
	Put(string, _Event) error
	Subscribe(string) (<-chan _Event, error)
	Shutdown()
}

var (
	_ _EventQueue = (*_InMemoryEventQueue)(nil)
)

func _MustNewEventQueue(topics []string) _EventQueue {
	return _MustNewInMemoryEventQueue(topics)
}
