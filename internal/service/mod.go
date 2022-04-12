package service

import (
	"context"

	"github.com/L-LYR/pns/internal/service/internal/dao"
)

func MustInitialize(ctx context.Context) {
	dao.MustInitTargetMongoDao(ctx)
	dao.MustInitLogRedisDao(ctx)
}

func MustShutdown(ctx context.Context) {
	dao.MustShutdownTargetMongoDao(ctx)
}
