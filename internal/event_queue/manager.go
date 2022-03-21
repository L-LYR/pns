package event_queue

type _EventQueueManager struct {
	Queue   EventQueue
	Workers []Worker
}

var (
	_Manager = &_EventQueueManager{
		Queue: NewEventQueue(
			[]string{
				_TargetEventTopic,
				_PushEventTopic,
			},
		),
		Workers: []Worker{
			_TargetEventWorker,
		},
	}
)

func MustInit() {
	_Manager.Queue.Start()
	for _, w := range _Manager.Workers {
		ch, err := _Manager.Queue.Subscribe(w.Topic())
		if err != nil {
			panic(err)
		}
		if err := w.RunOn(ch); err != nil {
			panic(err)
		}
	}
}

func MustShutdown() {
	_Manager.Queue.Shutdown()
	for _, w := range _Manager.Workers {
		w.Shutdown()
	}
}
