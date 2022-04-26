package model

import (
	"strconv"
	"strings"
)

type BizRuleStatus int8

const (
	BizRuleDisable BizRuleStatus = 0
	BizRuleEnable  BizRuleStatus = 1
)

type BizRule struct {
	Name        string
	Description string
	Salience    int
	Content     string
	Status      BizRuleStatus
}

func (r *BizRule) String() string {
	sb := &strings.Builder{}
	sb.WriteString("rule \"")
	sb.WriteString(r.Name)
	sb.WriteString("\" \"")
	sb.WriteString(r.Description)
	sb.WriteString("\" salience ")
	sb.WriteString(strconv.FormatInt(int64(r.Salience), 10))
	sb.WriteString("\nbegin\n")
	sb.WriteString(r.Content)
	sb.WriteString("\nend")
	return sb.String()
}
