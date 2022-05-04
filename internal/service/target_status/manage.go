package target_status

import (
	"context"
	"fmt"

	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/spaolacci/murmur3"
)

func _StatusOffset(key string) uint32 {
	return murmur3.Sum32([]byte(key))
}

type TargetStatusType int64

const (
	TargetOnline  TargetStatusType = 1
	TargetOffline TargetStatusType = 0

	_StatusStoreKey = "Status"
)

func CheckTargetOnline(ctx context.Context, appId int, deviceId string) bool {
	if status, err := dao.LogRedisDao.GetStatusBitmap(
		ctx,
		_StatusStoreKey,
		int64(_StatusOffset(fmt.Sprintf("%d:%s", appId, deviceId))),
	); err != nil {
		// by default, we think the target is online
		return true
	} else {
		return status == int64(TargetOnline)
	}
}

func SetTargetOnline(ctx context.Context, appId, deviceId string) error {
	return _SetTargetStatus(ctx, appId, deviceId, int(TargetOnline))
}

func SetTargetOffline(ctx context.Context, appId, deviceId string) error {
	return _SetTargetStatus(ctx, appId, deviceId, int(TargetOffline))
}

func _SetTargetStatus(ctx context.Context, appId, deviceId string, status int) error {
	return dao.LogRedisDao.SetStatusBitmap(
		ctx,
		_StatusStoreKey,
		int64(_StatusOffset(appId+":"+deviceId)),
		status,
	)
}
