package model

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type PushTaskType int8

const (
	DirectPush    PushTaskType = 1
	BroadcastPush PushTaskType = 2
	RangePush     PushTaskType = 3
)

func ParsePushTaskType(name string) (PushTaskType, error) {
	switch name {
	case "direct":
		return DirectPush, nil
	case "broadcast":
		return BroadcastPush, nil
	case "range":
		return RangePush, nil
	default:
		return 0, errors.New("unknown type")
	}
}

func (t PushTaskType) TopicNamePrefix() string {
	switch t {
	case DirectPush, RangePush:
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
	case RangePush:
		return "Range"
	default:
		panic("unreachable")
	}
}

type PushTaskStatusType int8

const (
	Pending  PushTaskStatusType = 0
	OnHandle PushTaskStatusType = 1
	Retry    PushTaskStatusType = 2
	Success  PushTaskStatusType = 3
	Failure  PushTaskStatusType = 4
	Filtered PushTaskStatusType = 5
)

type Qos = byte

const (
	AtMostOnce  Qos = 0
	AtLeastOnce Qos = 1
	ExactlyOnce Qos = 2
)

func ParseQos(s string) Qos {
	switch s {
	case "atMostOnce":
		return AtMostOnce
	case "atLeastOnce":
		return AtLeastOnce
	case "exactlyOnce":
		return ExactlyOnce
	default:
		return AtLeastOnce
	}
}

type PushTask interface {
	GetID() int64
	GetType() PushTaskType
	GetAppId() int
	GetPusher() PusherType
	GetMessage() *Message
	GetLogMeta() *LogMeta
	GetMeta() *PushTaskMeta
	GetStatus() PushTaskStatusType
	GetQos() Qos
	GetTopic() string

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

func (c *RetryCounter) SetRetry() {
	if c.Counter == NeverRetry {
		return
	}
	c.Counter--
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
	ID                int64              `json:"id"`
	Pusher            PusherType         `json:"pusher"`
	Qos               Qos                `json:"qos"`
	Status            PushTaskStatusType `json:"status"`
	CreationTime      time.Time          `json:"creationTime"`
	ValidationTime    time.Time          `json:"validationTime"`
	HandleTime        time.Time          `json:"handleTime"`
	EndTime           time.Time          `json:"endTime"`
	IgnoreFreqCtrl    bool               `json:"freqCtrl"`
	IgnoreOnlineCheck bool               `json:"onlineCheck"`
	*RetryCounter
}

func (m *PushTaskMeta) Spawn() *PushTaskMeta {
	return &PushTaskMeta{
		ID:             m.ID,
		Pusher:         m.Pusher,
		Qos:            m.Qos,
		Status:         m.Status,
		CreationTime:   time.Now(),
		IgnoreFreqCtrl: m.IgnoreFreqCtrl,
		RetryCounter: &RetryCounter{
			Counter: m.RetryCounter.Counter,
		},
	}
}

func (m *PushTaskMeta) SetRetry() {
	m.RetryCounter.SetRetry()
	m.Status = Retry
}
func (m *PushTaskMeta) SetSuccess()  { m.Status = Success }
func (m *PushTaskMeta) SetFailure()  { m.Status = Failure }
func (m *PushTaskMeta) SetOnHandle() { m.Status = OnHandle }
func (m *PushTaskMeta) SetPending()  { m.Status = Pending }
func (m *PushTaskMeta) SetFiltered() { m.Status = Filtered }

func (m *PushTaskMeta) UnderFreqCtrl() bool           { return !m.IgnoreFreqCtrl }
func (m *PushTaskMeta) NeedOnlineCheck() bool         { return !m.IgnoreOnlineCheck }
func (m *PushTaskMeta) IsRetry() bool                 { return m.Status == Retry }
func (m *PushTaskMeta) OnHandle() bool                { return m.Status == OnHandle }
func (m *PushTaskMeta) IsDone() bool                  { return m.Success() || m.Failure() }
func (m *PushTaskMeta) Success() bool                 { return m.Status == Success }
func (m *PushTaskMeta) Failure() bool                 { return m.Status == Failure }
func (m *PushTaskMeta) GetStatus() PushTaskStatusType { return m.Status }

func (m *PushTaskMeta) SetCreationTime(t time.Time)   { m.CreationTime = t }
func (m *PushTaskMeta) GetCreationTime() time.Time    { return m.CreationTime }
func (m *PushTaskMeta) SetValidationTime(t time.Time) { m.ValidationTime = t }
func (m *PushTaskMeta) GetValidationTime() time.Time  { return m.ValidationTime }
func (m *PushTaskMeta) SetHandleTime(t time.Time)     { m.HandleTime = t }
func (m *PushTaskMeta) GetHandleTime() time.Time      { return m.HandleTime }
func (m *PushTaskMeta) SetEndTime(t time.Time)        { m.EndTime = t }
func (m *PushTaskMeta) GetEndTime() time.Time         { return m.EndTime }
func (m *PushTaskMeta) GetID() int64                  { return m.ID }
func (m *PushTaskMeta) GetQos() Qos                   { return m.Qos }
func (m *PushTaskMeta) GetPusher() PusherType         { return m.Pusher }

// check type before use this
func AsDirectPushTask(t PushTask) *DirectPushTask {
	return t.(*DirectPushTask)
}

type DirectPushTask struct {
	*PushTaskMeta
	*Target
	*Message
}

func (t *DirectPushTask) GetType() PushTaskType  { return DirectPush }
func (t *DirectPushTask) GetAppId() int          { return t.App.ID }
func (t *DirectPushTask) GetMessage() *Message   { return t.Message }
func (t *DirectPushTask) GetMeta() *PushTaskMeta { return t.PushTaskMeta }
func (t *DirectPushTask) GetTopic() string {
	return fmt.Sprintf(
		"%s/%d/%s",
		t.GetType().TopicNamePrefix(),
		t.GetAppId(),
		t.Device.ID,
	)
}
func (t *DirectPushTask) GetLogMeta() *LogMeta {
	meta := &LogMeta{
		Type:     DirectPush,
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
	AppId int `json:"appId"`
	*PushTaskMeta
	*Message
}

func (t *BroadcastTask) GetType() PushTaskType  { return BroadcastPush }
func (t *BroadcastTask) GetAppId() int          { return t.AppId }
func (t *BroadcastTask) GetMessage() *Message   { return t.Message }
func (t *BroadcastTask) GetMeta() *PushTaskMeta { return t.PushTaskMeta }
func (t *BroadcastTask) GetTopic() string {
	return fmt.Sprintf(
		"%s/%d",
		t.GetType().TopicNamePrefix(),
		t.GetAppId(),
	)
}
func (t *BroadcastTask) GetLogMeta() *LogMeta {
	meta := &LogMeta{
		Type:   BroadcastPush,
		TaskId: t.ID,
		AppId:  t.AppId,
	}
	return meta
}

// check type before use this
func AsRangePushTask(t PushTask) *RangePushTask {
	return t.(*RangePushTask)
}

type RangePushTask struct {
	AppId         int `json:"appId"` // shared
	*PushTaskMeta     // owned
	*Message          // shared
	*FilterParams     // shared
	*Target           // owned
}

func (t *RangePushTask) GetType() PushTaskType  { return RangePush }
func (t *RangePushTask) GetAppId() int          { return t.AppId }
func (t *RangePushTask) GetMessage() *Message   { return t.Message }
func (t *RangePushTask) GetMeta() *PushTaskMeta { return t.PushTaskMeta }
func (t *RangePushTask) GetTopic() string {
	return fmt.Sprintf(
		"%s/%d/%s",
		t.GetType().TopicNamePrefix(),
		t.GetAppId(),
		t.Device.ID,
	)
}
func (t *RangePushTask) GetLogMeta() *LogMeta {
	meta := &LogMeta{
		Type:   RangePush,
		TaskId: t.ID,
		AppId:  t.AppId,
	}
	if t.Target != nil {
		meta.DeviceId = t.Target.Device.ID
	}
	return meta
}

func (t *RangePushTask) Spawn() *RangePushTask {
	return &RangePushTask{
		AppId:        t.AppId,
		PushTaskMeta: t.PushTaskMeta.Spawn(),
		Message:      t.Message,
		FilterParams: t.FilterParams,
		Target:       nil,
	}
}

type FilterParams struct {
	MinAppVersion *string   `json:"minAppVersion,omitempty"`
	MaxAppVersion *string   `json:"maxAppVersion,omitempty"`
	OsLimit       *[]string `json:"osLimit,omitempty"`
	BrandLimit    *[]string `json:"brandLimit,omitempty"`
}

// NOTICE:
// This function try to convert some fields into mongo condition
// TODO:
// refactor this kind of generator, it is not extensible.
// I think we shall change another way to convert parameters into conditions.
func (fp *FilterParams) CursorFilters() []*bson.E {
	filters := make([]*bson.E, 0)
	if fp.OsLimit != nil {
		filters = append(filters, &bson.E{
			Key: "os",
			Value: bson.M{
				"$in": fp.BrandLimit,
			},
		})
	}
	if fp.BrandLimit != nil {
		filters = append(filters, &bson.E{
			Key: "brand",
			Value: bson.M{
				"$in": fp.BrandLimit,
			},
		})
	}
	return filters
}

type PushTaskStage string

const (
	TaskCreation   PushTaskStage = "creation"
	TaskValidation PushTaskStage = "validation"
	TaskHandle     PushTaskStage = "handle"
	TaskDone       PushTaskStage = "done"
	TaskRetry      PushTaskStage = "retry"
	TaskRecv       PushTaskStage = "recv"
	TaskShow       PushTaskStage = "show"
)
