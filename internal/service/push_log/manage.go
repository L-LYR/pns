package log

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
)

// Task Request Log, sync
func PutTaskRequestLog(ctx context.Context, meta *model.PushLogMeta) error {
	return dao.LogRedisDao.AppendTaskEntry(ctx, meta)
}

// Generic Task Log, async
func PutTaskLog(ctx context.Context, l *model.LogEntry) error {
	return dao.LogRedisDao.AppendTaskLog(ctx, l)
}

// Get Task Log by Task ID
func GetTaskLogByID(ctx context.Context, id int) ([]*model.LogEntry, error) {
	return dao.LogRedisDao.GetTaskLogByID(ctx, id)
}

func GetTaskStatusByID(ctx context.Context, id int) (*model.LogEntry, error) {
	return dao.LogRedisDao.GetTaskLastLogByID(ctx, id)
}