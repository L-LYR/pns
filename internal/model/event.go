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

type LogEvent struct {
	Ctx   context.Context
	Entry *LogEntry
}

func (e *LogEvent) GetCtx() context.Context { return e.Ctx }
func (e *LogEvent) GetEntry() *LogEntry     { return e.Entry }
