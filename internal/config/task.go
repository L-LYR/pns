package config

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

var (
	_DefaultTaskQos model.Qos
)

func GetDefaultTaskQos() model.Qos {
	return _DefaultTaskQos
}

func MustLoadTaskDefaultConfig(ctx context.Context) {
	MustLoadDefaultTaskQos(ctx)
}

func MustLoadDefaultTaskQos(ctx context.Context) {
	var err error
	_DefaultTaskQos, err = model.ParseQos(MustLoadConfigValue(ctx, "task.default_qos").String())
	if err != nil {
		panic(err)
	}
}
