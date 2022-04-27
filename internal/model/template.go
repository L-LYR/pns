package model

import "regexp"

type MsgTpl struct {
	ID       int64                 `json:"templateId"`
	AppId    int                   `json:"appId"`
	Template map[string]*MsgTplStr `json:"template"`
}

type MsgTplStr struct {
	Str    string   `json:"s"`
	Params []string `json:"p"`
}

func NewMsgTplStr(s string) *MsgTplStr {
	// TODO: extract params
	result := &MsgTplStr{Str: s}
	for _, param := range _PlaceholderRegexp.FindAllStringSubmatch(s, -1) {
		result.Params = append(result.Params, param[1])
	}
	return result
}

var (
	// fields in message that can be template
	_TemplateFields = []string{
		"title",
		"content",
	}
	_PlaceholderRegexp = regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)
)

func FilterTemplateFields(tpl map[string]string) map[string]*MsgTplStr {
	filtered := make(map[string]*MsgTplStr)
	for _, field := range _TemplateFields {
		if v, ok := tpl[field]; ok {
			filtered[field] = NewMsgTplStr(v)
		}
	}
	return filtered
}
