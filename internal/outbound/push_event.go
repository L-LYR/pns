package outbound

import (
	"errors"

	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/model"
)

func PushEventConsumer(e event_queue.Event) error {
	pe, ok := e.(*model.PushEvent)
	if !ok {
		return errors.New("not PushEvent")
	}
	return PusherManager.Handle(pe.GetCtx(), pe.GetTask(), pe.PusherType())
}
