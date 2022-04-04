package util

import (
	"errors"
	"strconv"
	"strings"
)

const (
	_PusherNamePrefix = "pns-pusher"
	_TargetNamePrefix = "pns-target"
	_Seperator        = ":"
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

func GenerateTargetClientID(deviceId string, appId int) string {
	return strings.Join(
		[]string{
			_TargetNamePrefix,
			deviceId,
			strconv.FormatInt(int64(appId), 10),
		},
		_Seperator,
	)
}

func ParseClientID(src string) (appId int, isPusher bool, err error) {
	ss := strings.Split(src, _Seperator)
	if ss[0] == _PusherNamePrefix && len(ss) == 2 {
		if appId, err := strconv.Atoi(ss[1]); err != nil {
			return 0, false, err
		} else {
			return appId, true, nil
		}
	} else if ss[0] == _TargetNamePrefix && len(ss) == 3 {
		if appId, err := strconv.Atoi(ss[2]); err != nil {
			return 0, false, err
		} else {
			return appId, false, nil
		}
	}
	return 0, false, errors.New("invalid client id format")
}
