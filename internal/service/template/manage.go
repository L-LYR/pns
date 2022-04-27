package template

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
)

func CreateMessageTemplate(
	ctx context.Context,
	appId int,
	templateId int64,
	tpl map[string]*model.MsgTplStr,
) error {
	return dao.CreateMessageTemplate(ctx, &model.MsgTpl{
		ID:       templateId,
		AppId:    appId,
		Template: tpl,
	})
}

func QueryMessageTemplate(
	ctx context.Context, templateId int64,
) (*model.MsgTpl, error) {
	return dao.QueryMessageTemplate(ctx, templateId)
}
