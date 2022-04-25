package event_queue

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/monitor"
)

type _InMemoryEventQueue struct {
	cancellor context.CancelFunc
	working   bool
	cs        map[string]chan Event
}

func _MustNewInMemoryEventQueue(configs map[string]*config.EventQueueConfig) *_InMemoryEventQueue {
	q := &_InMemoryEventQueue{
		working: false,
		cs:      make(map[string]chan Event),
	}
	for _, cfg := range configs {
		q.addChannel(cfg)
	}
	return q
}

func (q *_InMemoryEventQueue) addChannel(c *config.EventQueueConfig) error {
	if _, ok := q.cs[c.Topic]; ok {
		return errors.New("duplicate queue")
	}
	q.cs[c.Topic] = make(chan Event, c.Length)
	return nil
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
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for topic, ch := range q.cs {
				monitor.EventQueueLength.WithLabelValues(topic).Set(float64(len(ch)))
			}
		}
	}
}
