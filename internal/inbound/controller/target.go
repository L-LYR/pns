package controller

import (
	"context"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/event_queue"
	v1 "github.com/L-LYR/pns/internal/inbound/api/v1"
	"github.com/L-LYR/pns/internal/local_storage"
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
	return nil, _UpsertTarget(ctx, req, model.CreateTarget)
}

func (api *_TargetAPI) UpdateTarget(
	ctx context.Context,
	req *v1.TargetUpdateReq,
) (*v1.TargetUpdateRes, error) {
	return nil, _UpsertTarget(ctx, req, model.UpdateTarget)
}

// NOTICE: request is *v1.TargetCreateReq or *v1.TargetUpdateReq
func _UpsertTarget(
	ctx context.Context,
	request interface{},
	t model.PushEventType,
) error {
	deviceInfo, appInfo, err := _ExtractInfos(ctx, request)
	if err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		return util.FinalError(gcode.CodeInvalidParameter, err, "Fail to extract infos")
	}
	target := &model.Target{Device: deviceInfo, App: appInfo}
	if err := event_queue.EventQueueManager.Put(
		config.TargetEventTopic(),
		&model.TargetEvent{Type: t, Ctx: ctx, Target: target},
	); err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		return util.FinalError(gcode.CodeInternalError, err, "Fail to emit target event")
	}
	return nil
}

func _ExtractInfos(ctx context.Context, request interface{}) (*model.Device, *model.App, error) {
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
	appName, ok := local_storage.GetAppNameByAppId(req.AppId)
	if !ok {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Unknown app id")
	}
	target, err := target.Query(ctx, req.DeviceId, appName)
	if err != nil {
		util.GLog.Errorf(ctx, "%v", err.Error())
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to query target info")
	}
	return &v1.TargetQueryRes{Target: target}, nil
}

func (api *_TargetAPI) GetToken(ctx context.Context, req *v1.TargetTokenReq) (*v1.TargetTokenRes, error) {
	// TODO：add expire time logic
	tokenSource := &util.TokenSource{}
	if err := copier.Copy(tokenSource, req); err != nil {
		return nil, err
	}
	token, err := util.NewTokenBuilder().Build(tokenSource)
	if err != nil {
		return nil, err
	}
	return &v1.TargetTokenRes{Token: token}, nil
}
