package event_queue

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/monitor"
)

const (
	_DefaultInMemoryChannelLength = 1000
)

type _InMemoryEventQueue struct {
	cancellor context.CancelFunc
	working   bool
	cs        map[string]chan Event
}

// TODO: add channel length monitor
func _MustNewInMemoryEventQueue(topics []string) *_InMemoryEventQueue {
	q := &_InMemoryEventQueue{
		working: false,
		cs:      make(map[string]chan Event),
	}
	for _, topic := range topics {
		q.cs[topic] = make(chan Event, _DefaultInMemoryChannelLength)
	}
	return q
}

func (q *_InMemoryEventQueue) Start(ctx context.Context) {
	q.working = true
	ctx, q.cancellor = context.WithCancel(ctx)
	go q.Monitor(ctx)
}

func (q *_InMemoryEventQueue) Put(topic string, e Event) error {
	if !q.working {
		return errors.New("event queue is down")
	}
	c, ok := q.cs[topic]
	if !ok {
		return errors.New("unknown topic")
	}
	c <- e
	return nil
}

func (q *_InMemoryEventQueue) Subscribe(topic string) (<-chan Event, error) {
	ch, ok := q.cs[topic]
	if !ok {
		return nil, errors.New("unknown topic")
	}
	return ch, nil
}

func (q *_InMemoryEventQueue) Shutdown(ctx context.Context) {
	q.working = false
}

func (q *_InMemoryEventQueue) Monitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second * 5):
			for topic, ch := range q.cs {
				monitor.EventQueueLength.WithLabelValues(topic).Set(float64(len(ch)))
			}
		}
	}
}
