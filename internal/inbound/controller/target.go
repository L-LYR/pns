package controller

import (
	"context"
	"time"

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

func (api *_TargetAPI) UpsertTarget(
	ctx context.Context,
	req *v1.UpdateTargetReq,
) (*v1.UpdateTargetRes, error) {
	targetInfo, err := _ExtractInfos(ctx, req)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInvalidParameter, err, "Fail to extract target info")
	}
	util.GLog.Printf(ctx, "%+v %+v", targetInfo.App, targetInfo.Device)
	if err := target.Upsert(ctx, targetInfo); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to update target info")
	}
	return nil, nil
}

func _ExtractInfos(
	ctx context.Context,
	req *v1.UpdateTargetReq,
) (*model.Target, error) {
	deviceInfo := &model.Device{}
	if err := copier.Copy(deviceInfo, req); err != nil {
		return nil, err
	}
	appInfo := &model.App{}
	if err := copier.Copy(appInfo, req); err != nil {
		return nil, err
	}
	return &model.Target{
		Device: deviceInfo,
		App:    appInfo,
		Tokens: &model.TokenSet{}, // here must new an empty token set
	}, nil
}

func (api *_TargetAPI) QueryTarget(
	ctx context.Context,
	req *v1.QueryTargetReq,
) (*v1.QueryTargetRes, error) {
	appName, ok := local_storage.GetAppNameByAppId(req.AppId)
	if !ok {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Unknown app id")
	}
	target, err := target.Query(ctx, appName, req.DeviceId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to query target info")
	}
	return &v1.QueryTargetRes{Target: target}, nil
}

// TODO: make token expire time into config
const (
	_TokenExpireTime = 7 * 24 * 3600 * time.Second
)

func _CheckTokenExpire(last time.Time) bool {
	return last.Add(_TokenExpireTime).Before(time.Now())
}

func (api *_TargetAPI) GetToken(
	ctx context.Context,
	req *v1.GetTokenReq,
) (*v1.GetTokenRes, error) {
	appName, ok := local_storage.GetAppNameByAppId(req.AppId)
	if !ok {
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Unknown app id")
	}

	targetInfo, err := target.Query(ctx, appName, req.DeviceId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to check device id")
	}

	if targetInfo == nil {
		return nil, util.FinalError(gcode.CodeInvalidRequest, nil, "Unknown target")
	}

	if !_CheckTokenExpire(targetInfo.TokenUpdateTime) {
		return &v1.GetTokenRes{Token: targetInfo.Tokens.Self}, nil
	}

	token, err := _GenerateToken(req.AppId, req.DeviceId)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to generate token")
	}

	if err := target.UpdateToken(ctx, appName, req.DeviceId, "self", token); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to update token")
	}
	return &v1.GetTokenRes{Token: token}, nil
}

func _GenerateToken(appId int, deviceId string) (string, error) {
	token, err := util.NewTokenBuilder().Build(
		&util.TokenSource{
			AppId:    appId,
			DeviceId: deviceId,
		},
	)
	if err != nil {
		return "", err
	}
	return token, nil
}
