// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
)

// MessageTemplate is the golang structure for table message_template.
type MessageTemplate struct {
	AppId      int         `json:"appId"      ` // app id
	TemplateId int         `json:"templateId" ` // template id
	Template   *gjson.Json `json:"template"   ` // message template
}