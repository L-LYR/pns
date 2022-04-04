package event_queue

import (
	"github.com/L-LYR/pns/internal/config"
)

type _EventQueueManager struct {
	topics  []string
	queue   EventQueue
	workers []_Worker
}

var (
	EventQueueManager = &_EventQueueManager{}
)

func (m *_EventQueueManager) MustRegister(cfg *config.ConsumerConfig, consumer Consumer) {
	m.topics = append(m.topics, cfg.Topic)
	m.workers = append(m.workers, _MustNewWorker(cfg, consumer))
}

func (m *_EventQueueManager) Put(topic string, event Event) error {
	return m.queue.Put(topic, event)
}

func (m *_EventQueueManager) MustStart() {
	m.queue = _MustNewEventQueue(m.topics...)
	m.queue.Start()
	for _, w := range m.workers {
		ch, err := m.queue.Subscribe(w.Topic())
		if err != nil {
			panic(err)
		}
		if err := w.RunOn(ch); err != nil {
			panic(err)
		}
	}
}

func (m *_EventQueueManager) MustShutdown() {
	m.queue.Shutdown()
	for _, w := range m.workers {
		w.Shutdown()
	}
}
