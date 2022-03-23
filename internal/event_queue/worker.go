package event_queue

import (
	"errors"
	"sync"

	"github.com/L-LYR/pns/internal/util"
)

type Consumer func(_Event) error

type _Worker interface {
	Topic() string
	RunOn(<-chan _Event) error
	Shutdown()
}

var (
	_ _Worker = (*_RealWorker)(nil)
)

func _MustNewWorker(topic string, dispatchN uint, consumer Consumer) _Worker {
	return &_RealWorker{
		topic:     topic,
		dispatchN: dispatchN,
		closeChan: make(chan struct{}),
		fn:        consumer,
	}
}

type _RealWorker struct {
	topic     string
	dispatchN uint
	wg        sync.WaitGroup
	closeChan chan struct{}
	fn        Consumer
}

func (w *_RealWorker) Topic() string { return w.topic }

func (w *_RealWorker) RunOn(ch <-chan _Event) error {
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
						util.GLog.Errorf(e.GetCtx(), "%+v", err)
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

func (w *_RealWorker) Shutdown() {
	w.closeChan <- struct{}{}
	w.wg.Wait()
}
