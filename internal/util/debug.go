package util

import (
	"strings"

	"github.com/gogf/gf/v2/os/genv"
	jsoniter "github.com/json-iterator/go"
)

func Readable(v interface{}) string {
	result, _ := jsoniter.MarshalIndent(v, "", strings.Repeat(" ", 4))
	return string(result)
}

func DebugOn() bool {
	return genv.Get("DEBUG").Bool()
}