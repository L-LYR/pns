// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// AppConfig is the golang structure of table app_config for DAO operations like Where/Data.
type AppConfig struct {
	g.Meta `orm:"table:app_config, do:true"`
	Id     interface{} // app id
	Name   interface{} // app name
}