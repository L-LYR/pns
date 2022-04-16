package model

import "errors"

type PusherType int

const (
	MQTTPusher PusherType = 1
	FCMPusher  PusherType = 2
	APNsPusher PusherType = 3
)

func ParsePusherType(name string) (PusherType, error) {
	switch name {
	case "mqtt":
		return MQTTPusher, nil
	case "fcm":
		return FCMPusher, nil
	case "apns":
		return APNsPusher, nil
	}
	return 0, errors.New("unknown pusher type")
}
