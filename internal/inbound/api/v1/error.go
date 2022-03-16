package v1

const (
	Success             int64 = 0
	InvalidParameters   int64 = 10001
	InternalServerError int64 = 10002
)

func ErrorMessage(code int64) string {
	switch code {
	case Success:
		return "success"
	case InvalidParameters:
		return "invalid parameters"
	case InternalServerError:
		return "internal error"
	default:
		return ""
	}
}

func RespondWith(code int64, payload ...interface{}) *CommonResponse {
	resp := &CommonResponse{
		ErrorCode:    code,
		ErrorMessage: ErrorMessage(code),
	}
	if len(payload) == 1 {
		resp.Payload = payload[0]
	} else {
		resp.Payload = payload
	}
	return resp
}
