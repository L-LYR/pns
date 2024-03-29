package log

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
)

func PutTaskEntry(ctx context.Context, meta *model.LogMeta) error {
	return dao.LogRedisDao.AppendTaskEntry(ctx, meta)
}

func PutTaskLog(ctx context.Context, l *model.LogEntry) error {
	return dao.LogRedisDao.AppendTaskLog(ctx, l)
}

func PutPushLog(ctx context.Context, l *model.LogEntry) error {
	if err := dao.LogRedisDao.CheckAndAppendTaskEntry(ctx, l.Meta); err != nil {
		return err
	}
	if err := dao.LogRedisDao.AppendPushLog(ctx, l); err != nil {
		return err
	}
	return dao.LogRedisDao.IncrTaskCounter(ctx, l.Meta.TaskId, string(l.Where))
}

func GetTaskLogByID(ctx context.Context, id int64) ([]*model.LogEntry, error) {
	return dao.LogRedisDao.GetTaskLogByID(ctx, id)
}

func GetTaskStatusByID(ctx context.Context, id int64) (*model.LogEntry, error) {
	return dao.LogRedisDao.GetTaskLastLogByID(ctx, id)
}

func GetPushLogByMeta(ctx context.Context, meta *model.LogMeta) ([]*model.LogEntry, error) {
	return dao.LogRedisDao.GetPushLogByMeta(ctx, meta)
}

func GetTaskEntryListByMeta(ctx context.Context, meta *model.LogMeta) ([]string, error) {
	return dao.LogRedisDao.GetTaskEntryListByMeta(ctx, meta)
}

func GetTaskStatisticsByID(ctx context.Context, id int64) (string, error) {
	result, err := dao.LogRedisDao.GetTaskStatistics(ctx, id, "send", "receive", "show")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d sent, %d received, %d showed", result[0], result[1], result[2]), nil
}

func CountLogEntry(ctx context.Context, meta *model.LogMeta, duration time.Duration) (int64, error) {
	now := time.Now()
	end := strconv.FormatInt(now.UnixMilli(), 10)
	begin := strconv.FormatInt(now.Add(-duration).UnixMilli(), 10)
	return dao.LogRedisDao.CountLogEntry(ctx, meta.EntryKey(), begin, end)
}
