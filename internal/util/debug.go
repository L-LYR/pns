package util

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Readable(v interface{}) string {
	result, _ := jsoniter.MarshalIndent(v, "", strings.Repeat(" ", 4))
	return string(result)
}
