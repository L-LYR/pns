package v1

type CommonResponse struct {
	ErrorCode    int64       `json:"errorCode:"`
	ErrorMessage string      `json:"errorMessage"`
	Payload      interface{} `json:"payload"`
}
