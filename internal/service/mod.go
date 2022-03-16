package service

import (
	"context"

	"github.com/L-LYR/pns/internal/service/internal/dao"
)

func MustInit(ctx context.Context) {
	dao.MustInit(ctx)
}

func MustShutdown(ctx context.Context) {
	dao.MustShutdown(ctx)
}
