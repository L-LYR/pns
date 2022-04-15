package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
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
	//TODO: check app created
	t, err := model.ParsePusherType(req.ConfigName)
	if err != nil {
		return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse pusher type")
	}

	var config interface{}
	switch t {
	case model.FCMPusher:
		FCMConfig := &model.FCMConfig{}
		if err := copier.Copy(FCMConfig, req.Config); err != nil {
			return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse config")
		}
		config = FCMConfig
	case model.APNsPusher:
		APNsConfig := &model.APNsConfig{}
		if err := copier.Copy(APNsConfig, req.Config); err != nil {
			return nil, util.FinalError(gcode.CodeValidationFailed, err, "Fail to parse config")
		}
		config = APNsConfig
	case model.MQTTPusher:
		return nil, util.FinalError(gcode.CodeInvalidParameter, nil, "Wrong request")
	default:
		panic("unreachable")
	}

	if err := app.CreateConfig(ctx, req.AppId, t, config); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to insert new app config")
	}
	return &v1.CreateAppConfigRes{}, nil
}

func (api *_AppAPI) OpenMQTT(ctx context.Context, req *v1.OpenMQTTReq) (*v1.OpenMQTTRes, error) {
	//TODO: check app created
	config := _GenerateRandMQTTConfig()
	if err := app.CreateConfig(ctx, req.AppId, model.MQTTPusher, config); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to insert new mqtt config")
	}
	return &v1.OpenMQTTRes{MQTTConfig: config}, nil
}

// TODO: this function will return the whole config content,
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

func _GenerateRandMQTTConfig() *v1.MQTTConfig {
	return &v1.MQTTConfig{
		// TODO: make the length of key and secret configurable
		PusherKey:      util.RandString(32),
		PusherSecret:   util.RandString(32),
		ReceiverKey:    util.RandString(32),
		ReceiverSecret: util.RandString(32),
	}
}
