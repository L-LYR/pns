package event_manager

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

type manager interface {
	Start()
	Put(*Event) error
	Subscribe(string, func(*Event) error) error
	Shutdown()
}

func newManager(topic string) manager {
	return newInMemoryEventQueue(topic)
}

// TODO: support RMQ

type InMemoryEventQueue struct {
	topic     string
	working   bool
	c         chan *Event
	closeChan chan struct{}
	once      sync.Once
}

// TODO: add option
// TODO: add channel length monitor
func newInMemoryEventQueue(topic string) *InMemoryEventQueue {
	return &InMemoryEventQueue{
		topic:     topic,
		c:         make(chan *Event, 1000),
		closeChan: make(chan struct{}),
		working:   false,
	}
}

func (q *InMemoryEventQueue) Start() {
	q.working = true
}

func (q *InMemoryEventQueue) Put(e *Event) error {
	if q.working {
		q.c <- e
		return nil
	} else {
		return errors.New("collector is down")
	}
}

func (q *InMemoryEventQueue) Subscribe(topic string, fn func(*Event) error) error {
	if q.topic != topic {
		return fmt.Errorf("unmatched topic, want %s, got %s", q.topic, topic)
	}
	q.once.Do(func() {
		go func() {
			for {
				select {
				case <-q.closeChan:
					return
				case event := <-q.c:
					if err := fn(event); err != nil {
						// TODO: maybe we can log the error cases
						g.Log().Line().Errorf(event.Ctx, "%+v", err)
					}
				}
			}
		}()
	})
	return nil
}

func (q *InMemoryEventQueue) Shutdown() {
	q.working = false
	q.closeChan <- struct{}{}
}
