package event_manager

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

type TargetEventType = int8

const (
	CreateTarget TargetEventType = 1
	UpdateTarget TargetEventType = 2
)

type TargetEventPayload struct {
	TargetEventType
	*model.Target
}

type Event struct {
	Ctx     context.Context
	Payload interface{}
}
