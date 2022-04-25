package event_queue

import (
	"context"
	"errors"
	"sync"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/util"
)

type Consumer func(Event) error

type _Worker interface {
	Topic() string
	RunOn(context.Context, <-chan Event) error
	Shutdown()
}

var (
	_ _Worker = (*_RealWorker)(nil)
)

func _MustNewWorker(cfg *config.EventConsumerConfig, consumer Consumer) _Worker {
	return &_RealWorker{
		cfg: cfg,
		fn:  consumer,
	}
}

type _RealWorker struct {
	cancellor context.CancelFunc
	cfg       *config.EventConsumerConfig
	wg        sync.WaitGroup
	fn        Consumer
}

func (w *_RealWorker) Topic() string { return w.cfg.Topic }

func (w *_RealWorker) RunOn(ctx context.Context, ch <-chan Event) error {
	if w.cfg == nil || !w.cfg.Check() || w.fn == nil {
		return errors.New("worker is uninitialized")
	}
	ctx, w.cancellor = context.WithCancel(ctx)
	for i := uint(0); i < w.cfg.Dispatch; i++ {
		w.wg.Add(1)
		go func() {
			for {
				select {
				case e := <-ch:
					if err := w.fn(e); err != nil {
						util.GLog.Errorf(e.GetCtx(), "%+v", err)
					}
				case <-ctx.Done():
					w.wg.Done()
					return
				}
			}
		}()
	}
	return nil
}

func (w *_RealWorker) Shutdown() {
	if w.cancellor != nil {
		w.cancellor()
		w.wg.Wait()
	}
}
