package model

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
)

type MsgTpl struct {
	ID       int64                 `json:"templateId"`
	AppId    int                   `json:"appId"`
	Template map[string]*MsgTplStr `json:"template"`
}

func (tpl *MsgTpl) FillInParams(
	replaceStrings map[string]map[string]string,
) (*Message, error) {
	fieldsMap := map[string]string{}
	for _, field := range _TemplateFields {
		tplStr, ok := tpl.Template[field]
		if !ok {
			continue
		}
		replaceString, ok := replaceStrings[field]
		if !ok {
			return nil, fmt.Errorf("replace strings of field %s are not given", field)
		}
		result, err := tplStr.FillInParams(replaceString)
		if err != nil {
			return nil, err
		}
		fieldsMap[field] = result
	}
	message := &Message{}
	if err := gconv.Struct(fieldsMap, message); err != nil {
		return nil, err
	}
	return message, nil
}

type MsgTplStr struct {
	Source       string   `json:"s"`
	Placeholders []string `json:"p"`
}

func NewMsgTplStr(s string) *MsgTplStr {
	result := &MsgTplStr{Source: s}
	for _, param := range _PlaceholderRegexp.FindAllStringSubmatch(s, -1) {
		result.Placeholders = append(result.Placeholders, param[1])
	}
	sort.Slice(result.Placeholders,
		func(i, j int) bool {
			return result.Placeholders[i] < result.Placeholders[j]
		},
	)
	return result
}

func (s *MsgTplStr) FillInParams(replaceStrings map[string]string) (string, error) {
	if len(replaceStrings) != len(s.Placeholders) {
		return "", errors.New("mismatched placeholders and replace strings")
	}
	for _, p := range s.Placeholders {
		if _, ok := replaceStrings[p]; !ok {
			return "", fmt.Errorf("placeholder %s is not given", p)
		}
	}
	return _PlaceholderRegexp.ReplaceAllStringFunc(
		s.Source,
		func(s string) string {
			return replaceStrings[strings.Trim(s, "${}")]
		},
	), nil
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
