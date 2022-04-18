package model

import "context"

type PushTaskEvent struct {
	Ctx  context.Context
	Task PushTask
}

func (e *PushTaskEvent) GetCtx() context.Context { return e.Ctx }
func (e *PushTaskEvent) GetTask() PushTask       { return e.Task }

type PushLogEvent struct {
	Ctx   context.Context
	Entry *LogEntry
}

func (e *PushLogEvent) GetCtx() context.Context { return e.Ctx }
func (e *PushLogEvent) GetEntry() *LogEntry     { return e.Entry }
