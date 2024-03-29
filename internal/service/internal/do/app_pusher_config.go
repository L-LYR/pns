// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// AppPusherConfig is the golang structure of table app_pusher_config for DAO operations like Where/Data.
type AppPusherConfig struct {
	g.Meta   `orm:"table:app_pusher_config, do:true"`
	AppId    interface{} // app id
	PusherId interface{} // pusher id
	Config   *gjson.Json // app pusher config
}
