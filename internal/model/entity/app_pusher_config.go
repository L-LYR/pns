// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
)

// AppPusherConfig is the golang structure for table app_pusher_config.
type AppPusherConfig struct {
	AppId    int         `json:"appId"    ` // app id
	PusherId int         `json:"pusherId" ` // pusher id
	Config   *gjson.Json `json:"config"   ` // app pusher config
}
