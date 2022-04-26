package model

import "context"

type PushTaskEvent struct {
	Ctx  context.Context
	Task PushTask
}

func (e *PushTaskEvent) GetCtx() context.Context { return e.Ctx }
func (e *PushTaskEvent) GetTask() PushTask       { return e.Task }

type LogEventType int8

const (
	PushLog LogEventType = 1
	TaskLog LogEventType = 2
)

type LogEvent struct {
	Ctx   context.Context
	Type  LogEventType
	Entry *LogEntry
}

func (e *LogEvent) GetCtx() context.Context { return e.Ctx }
func (e *LogEvent) GetEntry() *LogEntry     { return e.Entry }
func (e *LogEvent) GetType() LogEventType   { return e.Type }
