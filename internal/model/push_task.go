package model

type PushTaskType int8

const (
	PersonalPush  PushTaskType = 1
	BroadcastPush PushTaskType = 2
)

type PushTask struct {
	ID   uint64 `json:"id"`
	Type PushTaskType
	*Target
	*Message
	// TODO: Parameters
}
