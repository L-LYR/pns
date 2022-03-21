package event_queue

import (
	"errors"
	"sync"

	"github.com/L-LYR/pns/internal/util"
)

type Consumer func(Event) error

type Worker interface {
	Topic() string
	RunOn(<-chan Event) error
	Shutdown()
}

func _NewWorker(topic string, dispatchN uint, consumer Consumer) Worker {
	return &_Worker{
		topic:     topic,
		dispatchN: dispatchN,
		closeChan: make(chan struct{}),
		fn:        consumer,
	}
}

type _Worker struct {
	topic     string
	dispatchN uint
	wg        sync.WaitGroup
	closeChan chan struct{}
	fn        Consumer
}

func (w *_Worker) Topic() string { return w.topic }

func (w *_Worker) RunOn(ch <-chan Event) error {
	if w.dispatchN == 0 || w.fn == nil {
		return errors.New("worker is uninitialized")
	}
	for i := uint(0); i < w.dispatchN; i++ {
		w.wg.Add(1)
		go func() {
			for {
				select {
				case e := <-ch:
					if err := w.fn(e); err != nil {
						util.GLog.Error(e.GetCtx(), "%+v", err)
					}
				case <-w.closeChan:
					w.wg.Done()
					return
				}
			}
		}()
	}
	return nil
}

func (w *_Worker) Shutdown() {
	w.closeChan <- struct{}{}
	w.wg.Wait()
}
