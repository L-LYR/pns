package v1

type CommonRes struct {
	ErrorCode    int64       `json:"errorCode:"`
	ErrorMessage string      `json:"errorMessage"`
	Payload      interface{} `json:"payload"`
}
