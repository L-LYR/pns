package util

import "fmt"

type LogHandler func(s string)

var (
	_GlobalLogHandler LogHandler = nil
)

func SetGlobalLogHandler(fn LogHandler) { _GlobalLogHandler = fn }

func Log(format string, vs ...interface{}) {
	if _GlobalLogHandler == nil {
		return
	}
	_GlobalLogHandler(fmt.Sprintf(format, vs...))
}
