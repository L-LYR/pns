package controller

import (
	"net/http"

	"github.com/L-LYR/pns/internal/event_queue"
	"github.com/L-LYR/pns/internal/inbound/api/v1"
	"github.com/L-LYR/pns/internal/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jinzhu/copier"
)

func UpsertTarget(r *ghttp.Request) {
	ctx := r.GetCtx()
	req := &v1.TargetUpsertReq{}
	res := &v1.TargetUpsertRes{}
	defer r.Response.WriteJson(res)

	if err := r.Parse(req); err != nil {
		g.Log().Line().Errorf(ctx, "%v", err.Error())
		res.CommonRes = v1.RespondWith(v1.InvalidParameters)
		return
	}

	deviceInfo := &model.Device{}
	if err := copier.Copy(deviceInfo, req); err != nil {
		g.Log().Line().Errorf(ctx, "%v", err.Error())
		res.CommonRes = v1.RespondWith(v1.InternalServerError)
		return
	}
	appInfo := &model.App{}
	if err := copier.Copy(appInfo, req); err != nil {
		g.Log().Line().Errorf(ctx, "%v", err.Error())
		res.CommonRes = v1.RespondWith(v1.InternalServerError)
		return
	}

	if err := event_queue.EmitTargetEvent(&event_queue.TargetEvent{
		Ctx:     ctx,
		Type:    emitTargetEventType(r.Method),
		Payload: &model.Target{Device: deviceInfo, App: appInfo},
	}); err != nil {
		g.Log().Line().Errorf(ctx, "%v", err.Error())
		res.CommonRes = v1.RespondWith(v1.InternalServerError)
		return
	}

	res.CommonRes = v1.RespondWith(v1.Success)
}

func emitTargetEventType(m string) event_queue.TargetEventType {
	switch m {
	case http.MethodPatch, http.MethodPut:
		return event_queue.UpdateTarget
	case http.MethodPost:
		return event_queue.CreateTarget
	default:
		panic("unreachable")
	}
}
