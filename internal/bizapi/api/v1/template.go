package v1

import "github.com/gogf/gf/v2/frame/g"

type CreateMessageTemplateReq struct {
	g.Meta `path:"/message/template" method:"post"`
	AppId  int               `json:"appId" dc:"appId" v:"required|app-exist#app id is required"`
	Fields map[string]string `json:"fields" dc:"template fields" v:"required#template fields are required"`
}

type CreateMessageTemplateRes struct {
	TemplateID int64 `json:"templateId" dc:"template id, used in push task"`
}
