package model

type PushTaskType int8

const (
	PersonalPush  PushTaskType = 1
	BroadcastPush PushTaskType = 2
)

func (t PushTaskType) Name() string {
	switch t {
	case PersonalPush:
		return "PPush"
	case BroadcastPush:
		return "BPush"
	default:
		panic("unreachable")
	}
}

type PushTask struct {
	ID   int          `json:"id"`
	Type PushTaskType `json:"type"`
	*Target
	*Message
	// TODO: Parameters
}

func (t *PushTask) LogMeta() *PushLogMeta {
	return &PushLogMeta{
		TaskId:   t.ID,
		AppId:    t.App.ID,
		DeviceId: t.Device.ID,
	}
}
