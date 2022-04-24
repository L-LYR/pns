package config

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

func LoadGlobalConfig() {
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("config.toml")
}

func MustLoadConfig(ctx context.Context, name string, pointer interface{}) {
	if !g.Cfg().Available(ctx) {
		panic("global config is not avaliable")
	}
	if err := g.Cfg().MustGet(ctx, name).Struct(pointer); err != nil {
		panic(err)
	}
}

func MustLoadConfigValue(ctx context.Context, name string) *gvar.Var {
	if !g.Cfg().Available(ctx) {
		panic("global config is not avaliable")
	}
	return g.Cfg().MustGet(ctx, name)
}
