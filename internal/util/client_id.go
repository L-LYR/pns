package util

import (
	"errors"
	"strconv"
	"strings"
)

const (
	_TargetNamePrefix = "pns-target"
	_Seperator        = ":"
)

func GetRootClientID() string {
	return "pns-pusher:root"
}

func GetRootClientPass() string {
	return "pns_root"
}

func GetRootClientUser() string {
	return "pns_root"
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

func ParseClientID(src string) (appId int, err error) {
	ss := strings.Split(src, _Seperator)
	if ss[0] == _TargetNamePrefix && len(ss) == 3 {
		if appId, err := strconv.Atoi(ss[2]); err != nil {
			return 0, err
		} else {
			return appId, nil
		}
	}
	return 0, errors.New("invalid client id format")
}
