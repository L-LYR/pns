package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/app"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/jinzhu/copier"
)

var App = _AppAPI{}

type _AppAPI struct{}

func (api *_AppAPI) CreateApp(ctx context.Context, req *v1.CreateAppReq) (*v1.CreateAppRes, error) {
	if err := app.Create(ctx, req.AppName, req.AppId); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to insert new app")
	}
	return &v1.CreateAppRes{}, nil
}

func (api *_AppAPI) CreateConfig(ctx context.Context, req *v1.CreateAppConfigReq) (*v1.CreateAppConfigRes, error) {
	t, err := model.ParsePusherType(req.ConfigName)
	if err != nil {
		return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse pusher type")
	}

	delete(req.Config, "appId") // for fear that if the request config has a field named appId

	config := model.NewEmptyPusherConfig(req.AppId, t)
	switch t {
	case model.FCMPusher:
		if err := copier.Copy(config, req.Config); err != nil {
			return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse config")
		}
	case model.APNsPusher:
		if err := copier.Copy(config, req.Config); err != nil {
			return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse config")
		}
	case model.MQTTPusher:
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Wrong request")
	default:
		panic("unreachable")
	}

	if err := app.CreateConfig(ctx, config); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to insert new app config")
	}
	return &v1.CreateAppConfigRes{}, nil
}

func (api *_AppAPI) OpenMQTT(ctx context.Context, req *v1.OpenMQTTReq) (*v1.OpenMQTTRes, error) {
	config := _GenerateRandMQTTConfig(req.AppId)
	res := &v1.OpenMQTTRes{MQTTConfig: &v1.MQTTConfig{}}
	if err := copier.Copy(res.MQTTConfig, config); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to combine response")
	}
	if err := app.CreateConfig(ctx, config); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to insert new mqtt config")
	}
	return res, nil
}

// NOTICE: this function will return the whole config content,
func (api *_AppAPI) QueryConfig(ctx context.Context, req *v1.QueryAppConfigReq) (*v1.QueryAppConfigRes, error) {
	t, err := model.ParsePusherType(req.ConfigName)
	if err != nil {
		return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse pusher type")
	}
	config, err := app.QueryConfig(ctx, req.AppId, t)
	if err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to insert new mqtt config")
	}
	return &v1.QueryAppConfigRes{Config: config}, nil
}

func _GenerateRandMQTTConfig(appId int) *model.MQTTConfig {
	return &model.MQTTConfig{
		ID:             appId,
		PusherKey:      util.RandString(config.AuthKeyLength()),
		PusherSecret:   util.RandString(config.AuthSecretLength()),
		ReceiverKey:    util.RandString(config.AuthKeyLength()),
		ReceiverSecret: util.RandString(config.AuthSecretLength()),
	}
}
