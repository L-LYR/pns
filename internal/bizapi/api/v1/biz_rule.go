package v1

import "github.com/gogf/gf/v2/frame/g"

type InsertBizRuleReq struct {
	g.Meta      `path:"/rule" method:"post"`
	Name        string `json:"name" dc:"rule name" v:"required#rule name is required"`
	Description string `json:"description" dc:"rule description" v:"required#rule description is required"`
	Salience    int    `json:"salience" dc:"rule salience" v:"required|between:1,254#rule salience must be specified|salience should be in [1, 254]"`
	Content     string `json:"content" dc:"rule content" v:"required#rule content is required"`
}

type InsertBizRuleRes struct{}

type EnableRuleReq struct {
	g.Meta `path:"/rule/enable" method:"put"`
	Name   string `json:"name" dc:"rule name" v:"required#rule name is required"`
}

type EnableRuleRes struct{}

type DisableRuleReq struct {
	g.Meta `path:"/rule/disable" method:"put"`
	Name   string `json:"name" dc:"rule name" v:"required#rule name is required"`
}

type DisableRuleRes struct{}
