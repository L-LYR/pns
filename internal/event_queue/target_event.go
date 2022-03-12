package event_queue

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

type TargetEventType = int8

const (
	CreateTarget TargetEventType = 1
	UpdateTarget TargetEventType = 2
)

type TargetEvent struct {
	Type    TargetEventType
	Ctx     context.Context
	Payload *model.Target
}

func (e *TargetEvent) GetCtx() context.Context    { return e.Ctx }
func (e *TargetEvent) GetPayload() interface{}    { return e.Payload }
func (e *TargetEvent) EventType() TargetEventType { return e.Type }
