package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/bizcore"
	"github.com/L-LYR/pns/internal/model"
	bizrule "github.com/L-LYR/pns/internal/service/biz_rule"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Biz = &_BizAPI{}

type _BizAPI struct{}

func (api *_BizAPI) InsertBizRule(ctx context.Context, req *v1.InsertBizRuleReq) (*v1.InsertBizRuleRes, error) {
	if err := bizrule.InsertBizRule(
		ctx, req.Name, req.Description,
		req.Salience, req.Content,
	); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to insert biz rule")
	}
	return &v1.InsertBizRuleRes{}, nil
}

func (api *_BizAPI) EnableRule(
	ctx context.Context,
	req *v1.EnableRuleReq,
) (*v1.EnableRuleRes, error) {
	if err := bizrule.ChangeRuleStatus(
		ctx, req.Name, model.BizRuleEnable,
		func(ctx context.Context, rule *model.BizRule) error {
			return bizcore.AddNewRule(ctx, rule)
		},
	); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to enable rule")
	}
	return &v1.EnableRuleRes{}, nil
}

func (api *_BizAPI) DisableRule(
	ctx context.Context,
	req *v1.DisableRuleReq,
) (*v1.DisableRuleRes, error) {
	if err := bizrule.ChangeRuleStatus(
		ctx, req.Name, model.BizRuleDisable,
		func(ctx context.Context, rule *model.BizRule) error {
			return bizcore.RemoveRule(ctx, rule.Name)
		},
	); err != nil {
		return nil, util.FinalError(gcode.CodeInternalError, err, "Fail to disable rule")
	}
	return &v1.DisableRuleRes{}, nil
}
