package model

import "context"

type PushEvent struct {
	Ctx    context.Context
	Pusher PusherType
	Task   *PushTask
}

func (e *PushEvent) GetCtx() context.Context { return e.Ctx }
func (e *PushEvent) GetTask() *PushTask      { return e.Task }
func (e *PushEvent) PusherType() PusherType  { return e.Pusher }
