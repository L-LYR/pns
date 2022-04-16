package service

import (
	"context"

	"github.com/L-LYR/pns/internal/service/cache"
	"github.com/L-LYR/pns/internal/service/internal/dao"
)

func MustInitialize(ctx context.Context) {
	dao.MustInitTargetMongoDao(ctx)
	dao.MustInitLogRedisDao(ctx)
	cache.Config.MustLoad(ctx)
}

func MustShutdown(ctx context.Context) {
	dao.MustShutdownTargetMongoDao(ctx)
}
