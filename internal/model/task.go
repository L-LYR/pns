package model

import (
	"errors"
	"time"
)

type PushTaskType int8

const (
	DirectPush    PushTaskType = 1
	BroadcastPush PushTaskType = 2
)

func ParsePushTaskType(name string) (PushTaskType, error) {
	switch name {
	case "direct":
		return DirectPush, nil
	case "broadcast":
		return BroadcastPush, nil
	default:
		return 0, errors.New("unknown type")
	}
}

func (t PushTaskType) TopicNamePrefix() string {
	switch t {
	case DirectPush:
		return "DPush"
	case BroadcastPush:
		return "BPush"
	default:
		panic("unreachable")
	}
}

func (t PushTaskType) Name() string {
	switch t {
	case DirectPush:
		return "Direct"
	case BroadcastPush:
		return "Broadcast"
	default:
		panic("unreachable")
	}
}

type PushTaskStatusType int8

const (
	OnHandle PushTaskStatusType = 1
	Retry    PushTaskStatusType = 2
	Success  PushTaskStatusType = 3
	Failure  PushTaskStatusType = 4
)

type PushTask interface {
	GetID() int
	GetType() PushTaskType
	GetAppId() int
	GetPusher() PusherType
	GetMessage() *Message
	GetLogMeta() *LogMeta
	GetMeta() *PushTaskMeta
	GetStatus() PushTaskStatusType

	CanRetry() bool
}

type RetryTimes int

const (
	AlwaysRetry RetryTimes = -1
	NeverRetry  RetryTimes = 0
)

type RetryCounter struct {
	Counter RetryTimes `json:"retryCounter"`
	// -1 : always
	//  0 : never
	//  n : left times
}

func (c *RetryCounter) CanRetry() bool {
	if c.Counter == AlwaysRetry {
		return true
	}
	if c.Counter > NeverRetry {
		return true
	}
	return false
}

type PushTaskMeta struct {
	*RetryCounter
	Status       PushTaskStatusType `json:"status"`
	CreationTime time.Time          `json:"creationTime"`
	HandleTime   time.Time          `json:"handleTime"`
	EndTime      time.Time          `json:"endTime"`
}

func NewTaskMeta() *PushTaskMeta { return &PushTaskMeta{} }

func (m *PushTaskMeta) SetRetry() {
	if m.RetryCounter.Counter == NeverRetry {
		return
	}
	m.RetryCounter.Counter--
	m.Status = Retry
}
func (m *PushTaskMeta) SetSuccess() { m.Status = Success }
func (m *PushTaskMeta) SetFailure() { m.Status = Failure }
func (m *PushTaskMeta) SetOnHandle() { m.Status = OnHandle }

func (m *PushTaskMeta) IsRetry() bool                 { return m.Status == Retry }
func (m *PushTaskMeta) OnHandle() bool                { return m.Status == OnHandle }
func (m *PushTaskMeta) IsDone() bool                  { return m.Success() || m.Failure() }
func (m *PushTaskMeta) Success() bool                 { return m.Status == Success }
func (m *PushTaskMeta) Failure() bool                 { return m.Status == Failure }
func (m *PushTaskMeta) GetStatus() PushTaskStatusType { return m.Status }

func (m *PushTaskMeta) SetCreationTime(t time.Time) { m.CreationTime = t }
func (m *PushTaskMeta) GetCreationTime() time.Time  { return m.CreationTime }
func (m *PushTaskMeta) SetHandleTime(t time.Time)   { m.HandleTime = t }
func (m *PushTaskMeta) GetHandleTime() time.Time    { return m.HandleTime }
func (m *PushTaskMeta) SetEndTime(t time.Time)      { m.EndTime = t }
func (m *PushTaskMeta) GetEndTime() time.Time       { return m.EndTime }

func (m *PushTaskMeta) TotalDuration() time.Duration      { return m.EndTime.Sub(m.CreationTime) }
func (m *PushTaskMeta) ValidationDuration() time.Duration { return m.HandleTime.Sub(m.CreationTime) }
func (m *PushTaskMeta) HandleDuration() time.Duration     { return m.EndTime.Sub(m.HandleTime) }

// check type before use this
func AsDirectPush(t PushTask) *DirectPushTask {
	return t.(*DirectPushTask)
}

type DirectPushTask struct {
	ID     int        `json:"id"`
	Pusher PusherType `json:"pusher"`
	*PushTaskMeta
	*Target
	*Message
}

func (t *DirectPushTask) GetID() int             { return t.ID }
func (t *DirectPushTask) GetType() PushTaskType  { return DirectPush }
func (t *DirectPushTask) GetAppId() int          { return t.App.ID }
func (t *DirectPushTask) GetPusher() PusherType  { return t.Pusher }
func (t *DirectPushTask) GetMessage() *Message   { return t.Message }
func (t *DirectPushTask) GetMeta() *PushTaskMeta { return t.PushTaskMeta }
func (t *DirectPushTask) GetLogMeta() *LogMeta {
	meta := &LogMeta{
		TaskId:   t.ID,
		AppId:    t.App.ID,
		DeviceId: t.Device.ID,
	}
	return meta
}

// check type before use this
func AsBroadcastTask(t PushTask) *BroadcastTask {
	return t.(*BroadcastTask)
}

type BroadcastTask struct {
	ID     int        `json:"id"`
	AppId  int        `json:"appId"`
	Pusher PusherType `json:"pusher"`
	*PushTaskMeta
	*Message
	// FilterParams
}

func (t *BroadcastTask) GetID() int             { return t.ID }
func (t *BroadcastTask) GetType() PushTaskType  { return BroadcastPush }
func (t *BroadcastTask) GetAppId() int          { return t.AppId }
func (t *BroadcastTask) GetPusher() PusherType  { return t.Pusher }
func (t *BroadcastTask) GetMessage() *Message   { return t.Message }
func (t *BroadcastTask) GetMeta() *PushTaskMeta { return t.PushTaskMeta }
func (t *BroadcastTask) GetLogMeta() *LogMeta {
	meta := &LogMeta{
		TaskId: t.ID,
		AppId:  t.AppId,
	}
	return meta
}
