package middleware

import (
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

type CommonResponse struct {
	ErrorCode    int         `json:"errorCode:" dc:"error code"`
	ErrorMessage string      `json:"errorMessage" dc:"error message"`
	Payload      interface{} `json:"payload" dc:"response payload"`
}

func RespondWith(code gcode.Code, payload ...interface{}) *CommonResponse {
	resp := &CommonResponse{}
	resp.ErrorCode = code.Code()
	resp.ErrorMessage = code.Message()
	if len(payload) == 1 {
		resp.Payload = payload[0]
	} else {
		resp.Payload = payload
	}
	return resp
}

func CommonResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()
	ctx := r.GetCtx()
	err := r.GetError()
	code := gerror.Code(err)
	if code == gcode.CodeNil && err != nil {
		code = gcode.CodeInternalError
	} else if err == nil {
		code = gcode.CodeOK
	}
	if err != nil {
		util.GLog.Errorf(ctx, "Final Error: %+v, Detail: %+v", err.Error(), code.Detail())
	}

	if r.Response.BufferLength() > 0 {
		// response is ready, do nothing here
		return
	}

	res := r.GetHandlerResponse()
	if err := r.Response.WriteJson(RespondWith(code, res)); err != nil {
		util.GLog.Warningf(ctx, "%+v", err.Error())
	}
}
