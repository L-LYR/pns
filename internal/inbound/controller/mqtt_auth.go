package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/inbound/api/v1"
	"github.com/L-LYR/pns/internal/service/app"
	"github.com/gogf/gf/v2/frame/g"
)

var MQTTAuth = &_MQTTAuthAPI{}

type _MQTTAuthAPI struct{}

func (api *_MQTTAuthAPI) Auth(
	ctx context.Context,
	req *v1.MQTTAuthReq,
) (*v1.MQTTAuthRes, error) {
	ok, reason := app.Authorization(ctx, req.Username, req.Password, req.ClientId)
	res := &v1.MQTTAuthRes{CommonAuthRes: v1.CommonAuthRes{Ok: ok, Error: reason}}
	if err := g.RequestFromCtx(ctx).Response.WriteJson(res); err != nil {
		return nil, err
	}
	return res, nil
}

var ACLCheck = &_ACLCheckAPI{}

type _ACLCheckAPI struct{}

func (api *_MQTTAuthAPI) Check(
	ctx context.Context,
	req *v1.ACLCheckReq,
) (*v1.ACLCheckRes, error) {
	ok, reason := app.ACLCheck(ctx, req.Username, req.ClientId)
	res := &v1.ACLCheckRes{CommonAuthRes: v1.CommonAuthRes{Ok: ok, Error: reason}}
	// ignore general middleware
	if err := g.RequestFromCtx(ctx).Response.WriteJson(res); err != nil {
		return nil, err
	}
	return res, nil
}
