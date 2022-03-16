package controller

import (
	"context"

	"github.com/L-LYR/pns/internal/event_queue"
	v1 "github.com/L-LYR/pns/internal/inbound/api/v1"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jinzhu/copier"
)

func CreateTarget(r *ghttp.Request) {
	ctx := r.GetCtx()
	req := &v1.TargetCreateRequest{}
	if err := r.Parse(req); err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InvalidParameters))
		return
	}
	deviceInfo, appInfo, err := extractInfos(ctx, req)
	if err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InvalidParameters))
		return
	}
	if err := emitTargetEvent(
		ctx,
		deviceInfo,
		appInfo,
		event_queue.CreateTarget,
	); err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InternalServerError))
		return
	}
	r.Response.WriteJson(v1.RespondWith(v1.Success))
}

func UpdateTarget(r *ghttp.Request) {
	ctx := r.GetCtx()
	req := &v1.TargetUpdateRequest{}
	if err := r.Parse(req); err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InvalidParameters))
		return
	}
	deviceInfo, appInfo, err := extractInfos(ctx, req)
	if err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InvalidParameters))
		return
	}
	if err := emitTargetEvent(
		ctx,
		deviceInfo,
		appInfo,
		event_queue.UpdateTarget,
	); err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InternalServerError))
		return
	}
	r.Response.WriteJson(v1.RespondWith(v1.Success))
}

// NOTICE: request is *v1.TargetCreateRequest or *v1.TargetUpdateRequest
func extractInfos(ctx context.Context, request interface{}) (*model.Device, *model.App, error) {
	deviceInfo := &model.Device{}
	if err := copier.Copy(deviceInfo, request); err != nil {
		return nil, nil, err
	}
	appInfo := &model.App{}
	if err := copier.Copy(appInfo, request); err != nil {
		return nil, nil, err
	}
	return deviceInfo, appInfo, nil
}

func emitTargetEvent(
	c context.Context,
	d *model.Device,
	a *model.App,
	t event_queue.TargetEventType,
) error {
	return event_queue.EmitTargetEvent(&event_queue.TargetEvent{
		Ctx:     c,
		Type:    t,
		Payload: &model.Target{Device: d, App: a},
	})
}

func QueryTarget(r *ghttp.Request) {
	ctx := r.GetCtx()
	req := &v1.TargetQueryRequest{}
	if err := r.ParseQuery(req); err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InvalidParameters))
		return
	}

	target, err := target.Query(ctx, req.DeviceId, req.AppId)
	if err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		r.Response.WriteJson(v1.RespondWith(v1.InternalServerError))
		return
	}

	r.Response.WriteJson(v1.RespondWith(v1.Success, target))
}
