package event_queue

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

type PushEventType = int8

type PushEvent struct {
	Type    PushEventType
	Ctx     context.Context
	Payload *model.PushTask
}

func (e *PushEvent) GetCtx() context.Context  { return e.Ctx }
func (e *PushEvent) GetPayload() interface{}  { return e.Payload }
func (e *PushEvent) EventType() PushEventType { return e.Type }
