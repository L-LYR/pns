package controller

import (
	"context"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/template"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
)

var Template = _TemplateAPI{}

type _TemplateAPI struct{}

func (api *_TemplateAPI) CreateMessageTemplate(
	ctx context.Context,
	req *v1.CreateMessageTemplateReq,
) (*v1.CreateMessageTemplateRes, error) {

	id := util.GenerateTemplateId()
	if err := template.CreateMessageTemplate(
		ctx, req.AppId, id,
		model.FilterTemplateFields(req.Fields),
	); err != nil {
		return nil, util.FinalError(
			gcode.CodeInternalError,
			err, "Fail to insert template",
		)
	}

	return &v1.CreateMessageTemplateRes{TemplateID: id}, nil
}
