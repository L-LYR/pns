package model

import (
	"errors"
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

type PushTask interface {
	GetID() int
	GetType() PushTaskType
	GetAppId() int
	GetPusher() PusherType
	GetMessage() *Message
	GetLogMeta() *LogMeta

	Retry() bool
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

func (c *RetryCounter) Retry() bool {
	if c.Counter == AlwaysRetry {
		return true
	}
	if c.Counter > NeverRetry {
		c.Counter--
		return true
	}
	return false
}

// check type before use this
func AsDirectPush(t PushTask) *DirectPushTask {
	return t.(*DirectPushTask)
}

type DirectPushTask struct {
	ID     int        `json:"id"`
	Pusher PusherType `json:"pusher"`
	*RetryCounter
	*Target
	*Message
}

func (t *DirectPushTask) GetID() int            { return t.ID }
func (t *DirectPushTask) GetType() PushTaskType { return DirectPush }
func (t *DirectPushTask) GetAppId() int         { return t.App.ID }
func (t *DirectPushTask) GetPusher() PusherType { return t.Pusher }
func (t *DirectPushTask) GetMessage() *Message  { return t.Message }
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
	*RetryCounter
	*Message
	// FilterParams
}

func (t *BroadcastTask) GetID() int            { return t.ID }
func (t *BroadcastTask) GetType() PushTaskType { return BroadcastPush }
func (t *BroadcastTask) GetAppId() int         { return t.AppId }
func (t *BroadcastTask) GetPusher() PusherType { return t.Pusher }
func (t *BroadcastTask) GetMessage() *Message  { return t.Message }
func (t *BroadcastTask) GetLogMeta() *LogMeta {
	meta := &LogMeta{
		TaskId: t.ID,
		AppId:  t.AppId,
	}
	return meta
}
