package model

type PushTask struct {
	ID uint64 `json:"id"`
	*Target
	*Message
	// TODO: Parameters
}
