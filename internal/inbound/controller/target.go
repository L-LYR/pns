package controller

import (
	"context"

	"github.com/L-LYR/pns/internal/event_queue"
	v1 "github.com/L-LYR/pns/internal/inbound/api/v1"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/jinzhu/copier"
)

var Target = &_TargetAPI{}

type _TargetAPI struct{}

func (api *_TargetAPI) CreateTarget(
	ctx context.Context,
	req *v1.TargetCreateReq,
) (*v1.TargetCreateRes, error) {
	return nil, upsertTarget(ctx, req, event_queue.CreateTarget)
}

func (api *_TargetAPI) UpdateTarget(
	ctx context.Context,
	req *v1.TargetUpdateReq,
) (*v1.TargetUpdateRes, error) {
	return nil, upsertTarget(ctx, req, event_queue.UpdateTarget)
}

// NOTICE: request is *v1.TargetCreateReq or *v1.TargetUpdateReq
func upsertTarget(
	ctx context.Context,
	request interface{},
	t event_queue.PushEventType,
) error {
	deviceInfo, appInfo, err := extractInfos(ctx, request)
	if err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		return util.FinalError(gcode.CodeInvalidParameter, err, "Fail to extract infos")
	}
	if err := event_queue.SendTargetEvent(
		ctx,
		&model.Target{Device: deviceInfo, App: appInfo},
		t,
	); err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		return util.FinalError(gcode.CodeInternalError, err, "Fail to emit target event")
	}
	return nil
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

func (api *_TargetAPI) QueryTarget(ctx context.Context, req *v1.TargetQueryReq) (*v1.TargetQueryRes, error) {
	target, err := target.Query(ctx, req.DeviceId, req.AppId)
	if err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to query target info")
	}
	return &v1.TargetQueryRes{Target: target}, nil
}
