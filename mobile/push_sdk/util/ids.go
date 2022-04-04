package util

import (
	"fmt"

	"github.com/google/uuid"
)

/*
	'device id' is an external id which is used to
	mark a unique device, here we use uuid to mock this.
*/

func GenerateDeviceId() string {
	return uuid.New().String()
}

func GenerateClientId(prefix string, deviceId string, appId int) string {
	return fmt.Sprintf("%s:%s:%d", prefix, deviceId, appId)
}
