package model

import "context"

type PushEventType = int8

const (
	Push PushEventType = 1
)

type PushEvent struct {
	Type   PushEventType
	Ctx    context.Context
	Pusher PusherType
	Task   *PushTask
}

func (e *PushEvent) GetCtx() context.Context  { return e.Ctx }
func (e *PushEvent) GetTask() *PushTask       { return e.Task }
func (e *PushEvent) EventType() PushEventType { return e.Type }
func (e *PushEvent) PusherType() PusherType   { return e.Pusher }

type TargetEventType = int8

const (
	CreateTarget TargetEventType = 1
	UpdateTarget TargetEventType = 2
)

type TargetEvent struct {
	Type   TargetEventType
	Ctx    context.Context
	Target *Target
}

func (e *TargetEvent) GetCtx() context.Context    { return e.Ctx }
func (e *TargetEvent) GetTarget() *Target         { return e.Target }
func (e *TargetEvent) EventType() TargetEventType { return e.Type }
