package model

type PushTask struct {
	ID string `json:"id"`
	*Target
	*Message
	// TODO: Parameters
}
