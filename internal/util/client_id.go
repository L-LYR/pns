package util

import (
	"errors"
	"strconv"
	"strings"
)

const (
	_PusherNamePrefix = "pns_pusher"
	_TargetNamePrefix = "pns_target"
	_Seperator        = "-"
)

func GeneratePusherClientID(appId int) string {
	return strings.Join(
		[]string{
			_PusherNamePrefix,
			strconv.FormatInt(int64(appId), 10),
		},
		_Seperator,
	)
}

func GenerateDeviceClientID(deviceId string, appId int) string {
	return strings.Join(
		[]string{
			_TargetNamePrefix,
			deviceId,
			strconv.FormatInt(int64(appId), 10),
		},
		_Seperator,
	)
}

func ParseClientID(src string) (int, error) {
	ss := strings.Split(src, _Seperator)
	if ss[0] == _PusherNamePrefix && len(ss) == 2 {
		return strconv.Atoi(ss[1])
	} else if ss[0] == _TargetNamePrefix && len(ss) == 3 {
		return strconv.Atoi(ss[2])
	}
	return 0, errors.New("invalid client id format")
}
