package event_queue

import (
	"context"

	"github.com/L-LYR/pns/internal/config"
)

type _EventQueueManager struct {
	queue   EventQueue
	configs map[string]*config.EventQueueConfig
	workers map[string]_Worker
}

var (
	EventQueueManager = &_EventQueueManager{
		configs: make(map[string]*config.EventQueueConfig),
		workers: make(map[string]_Worker),
	}
)

func (m *_EventQueueManager) MustRegister(cfg *config.EventQueueConfig, consumer Consumer) {
	m.configs[cfg.Topic] = cfg
	m.workers[cfg.Topic] = _MustNewWorker(cfg, consumer)
}

func (m *_EventQueueManager) Put(topic string, event Event) error {
	return m.queue.Put(topic, event)
}

func (m *_EventQueueManager) MustStart(ctx context.Context) {
	m.queue = _MustNewEventQueue(m.configs)

	m.queue.Start(ctx)
	for _, w := range m.workers {
		ch, err := m.queue.Subscribe(w.Topic())
		if err != nil {
			panic(err)
		}
		if err := w.RunOn(ctx, ch); err != nil {
			panic(err)
		}
	}
}

func (m *_EventQueueManager) MustShutdown(ctx context.Context) {
	m.queue.Shutdown(ctx)
	for _, w := range m.workers {
		w.Shutdown()
	}
}
