package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// TODO: simplify

type DirectPushBase struct {
	AppId             int    `json:"appId" dc:"registered app id" v:"required|app-exist#app id is required|unknown app"`
	DeviceId          string `json:"deviceId" dc:"available device id, only" v:"required#device id is required"`
	Retry             int    `json:"retry" dc:"retry times, -1 for always retry" v:"min:-1#retry time is greater than -1"`
	IgnoreFreqCtrl    bool   `json:"ignoreFreqCtrl" dc:"ignore frequency control"`
	IgnoreOnlineCheck bool   `json:"ignoreOnlineCheck" dc:"ignore online check"`
}

type BroadcastPushBase struct {
	AppId          int  `json:"appId" dc:"registered app id" v:"required|app-exist#app id is required"`
	IgnoreFreqCtrl bool `json:"ignoreFreqCtrl" dc:"ignore frequency control"`
}

type BasicMessage struct {
	Title   string `json:"title" dc:"push message title" v:"required#title is required"`
	Content string `json:"content" dc:"push message content" v:"required#content is required"`
}

type TemplateMessage struct {
	// MessageFieldsNotInTemplate
	Id         int64                   `json:"id" dc:"push message template id, if setted, we will use it as the message base" v:"required"`
	ParamLists map[string]ParamWrapper `json:"params" dc:"push message template parameters" v:"required"`
}

// TODO: this is because goframe have bug in gconv
// see https://github.com/gogf/gf/issues/1708
type ParamWrapper struct {
	PR map[string]string `json:"pr" dc:"placeholder to replace string" v:"required"`
}

type FilterParams struct {
	// TODO: add more
	MinAppVersion *string   `json:"minAppVersion,omitempty"`
	MaxAppVersion *string   `json:"maxAppVersion,omitempty"`
	OsLimit       *[]string `json:"osLimit,omitempty"`
	BrandLimit    *[]string `json:"brandLimit,omitempty"`
}

type PushResBase struct {
	PushTaskId int64 `json:"pushTaskId" dc:"use task id to tracing it"`
}

type DirectPushReq struct {
	g.Meta `path:"/push/direct" method:"post"`
	*DirectPushBase
	Message *BasicMessage `json:"message" dc:"push message" v:"required#message is required"`
}

type DirectPushRes struct{ *PushResBase }

type TemplateDirectPushReq struct {
	g.Meta `path:"/push/template/direct" method:"post"`
	*DirectPushBase
	Message *TemplateMessage `json:"message" dc:"push template message" v:"required#message is required"`
}

type TemplateDirectPushRes struct{ *PushResBase }

type BroadcastPushReq struct {
	g.Meta `path:"/push/broadcast" method:"post"`
	*BroadcastPushBase
	Message *BasicMessage `json:"message" dc:"push message" v:"required#message is required"`
}

type BroadcastPushRes struct{ *PushResBase }

type TemplateBroadcastPushReq struct {
	g.Meta `path:"/push/template/broadcast" method:"post"`
	*BroadcastPushBase
	Message *TemplateMessage `json:"message" dc:"push message" v:"required#message is required"`
}

type TemplateBroadcastPushRes struct{ *PushResBase }

type RangePushReq struct {
	g.Meta `path:"/push/range" method:"post"`
	*BroadcastPushBase
	*FilterParams `json:"filterParams" dc:"push filter parameters"`
	Message       *BasicMessage `json:"message" dc:"push message" v:"required#message is required"`
}

type RangePushRes struct{ *PushResBase }

type TemplateRangePushReq struct {
	g.Meta `path:"/push/template/range" method:"post"`
	*BroadcastPushBase
	*FilterParams `json:"filterParams" dc:"push filter parameters"`
	Message       *TemplateMessage `json:"message" dc:"push message" v:"required#message is required"`
}

type TemplateRangePushRes struct{ *PushResBase }
