package util

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func FinalError(baseCode gcode.Code, detail interface{}, explanations ...string) error {
	return gerror.NewCode(gcode.WithCode(baseCode, detail), explanations...)
}
