package internal

import (
	"strings"

	"github.com/L-LYR/pns/internal/model"
)

/*
	Functions exposed to rules should be wrapperred or defined here
	and be registered in _ApiOuter.
	These functions are not exposed to any other packages, but only used in rules.
*/

func OuterApis() map[string]interface{} {
	return _ApiOuter
}

func PredefinedRules() string {
	return strings.Join(_PredefinedRules, "")
}

var (
	_ApiOuter = map[string]interface{}{
		"BeginObserve":     _BeginObserve,
		"EndObserve":       _EndObserve,
		"FrequencyControl": _FreqCtrl,
		"AsDirectPush":     model.AsDirectPushTask,
		"AsBroadcastPush":  model.AsBroadcastTask,
		"AsRangePush":      model.AsRangePushTask,
	}

	// sorted by salience
	_PredefinedRules = []string{
		_BeginRule,
		_EndRule,
		_FreqCtrlRule,
	}
)
