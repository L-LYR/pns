package internal

import (
	"context"
	"errors"
	"sync"

	"github.com/L-LYR/pns/internal/util"
)

type Event interface {
	GetCtx() context.Context
	GetPayload() interface{}
}

const (
	_DefaultInMemoryChannelLength = 1000
)

type EventQueue interface {
	Start()
	Put(string, Event) error
	Subscribe(string, func(Event) error) error
	Shutdown()
}

func NewEventQueue(topics []string) EventQueue {
	return newInMemoryEventQueue(topics)
}

// TODO: support RMQ

type InMemoryEventQueue struct {
	working   bool
	cs        map[string]chan Event
	wgs       map[string]sync.WaitGroup
	closeChan chan struct{}
}

// TODO: add option
// TODO: add channel length monitor
func newInMemoryEventQueue(topics []string) *InMemoryEventQueue {
	q := &InMemoryEventQueue{
		closeChan: make(chan struct{}),
		working:   false,
		cs:        make(map[string]chan Event),
		wgs:       make(map[string]sync.WaitGroup),
	}
	for _, topic := range topics {
		q.cs[topic] = make(chan Event, _DefaultInMemoryChannelLength)
		q.wgs[topic] = sync.WaitGroup{}
	}
	return q
}

func (q *InMemoryEventQueue) Start() {
	q.working = true
}

func (q *InMemoryEventQueue) Put(topic string, e Event) error {
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

func (q *InMemoryEventQueue) Subscribe(topic string, fn func(Event) error) error {
	wg, ok := q.wgs[topic]
	if !ok {
		return errors.New("unknown topic")
	}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-q.closeChan:
				wg.Done()
				return
			case event := <-q.cs[topic]:
				if err := fn(event); err != nil {
					// TODO: maybe we can log the error cases
					util.GLog.Errorf(event.GetCtx(), "%+v", err)
				}
			}
		}
	}()
	return nil
}

func (q *InMemoryEventQueue) Shutdown() {
	q.working = false
	q.closeChan <- struct{}{}
	for _, wg := range q.wgs {
		wg.Wait()
	}
}
