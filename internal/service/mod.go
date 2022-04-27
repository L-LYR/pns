package service

import (
	"context"

	"github.com/L-LYR/pns/internal/service/cache"
	"github.com/L-LYR/pns/internal/service/internal/dao"
)

func MustInitialize(ctx context.Context) {
	dao.MustInitializeTargetMongoDao(ctx)
	dao.MustInitLogRedisDao(ctx)
	cache.Config.MustLoad(ctx)
	cache.MsgTpl.MustInitialize(ctx)
}

func MustShutdown(ctx context.Context) {
	dao.MustShutdownTargetMongoDao(ctx)
}
