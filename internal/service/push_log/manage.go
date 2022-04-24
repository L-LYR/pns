package log

import (
	"context"
	"fmt"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
)

func PutTaskEntry(ctx context.Context, meta *model.LogMeta) error {
	return dao.LogRedisDao.AppendTaskEntry(ctx, meta)
}

func PutTaskLog(
	ctx context.Context,
	meta *model.LogMeta,
	where string,
	hint string,
) error {
	return dao.LogRedisDao.AppendTaskLog(ctx, &model.LogEntry{
		LogBase: &model.LogBase{
			Meta:  meta,
			T:     time.Now().UnixMilli(),
			Where: where,
		},
		Hint: hint,
	})
}

// async
func PutPushLog(ctx context.Context, l *model.LogEntry) error {
	if err := dao.LogRedisDao.CheckAndAppendTaskEntry(ctx, l.Meta); err != nil {
		return err
	}
	return dao.LogRedisDao.AppendPushLog(ctx, l)
}

func GetTaskLogByID(ctx context.Context, id int) ([]*model.LogEntry, error) {
	return dao.LogRedisDao.GetTaskLogByID(ctx, id)
}

func GetTaskStatusByID(ctx context.Context, id int) (*model.LogEntry, error) {
	return dao.LogRedisDao.GetTaskLastLogByID(ctx, id)
}

func GetPushLogByMeta(ctx context.Context, meta *model.LogMeta) ([]*model.LogEntry, error) {
	return dao.LogRedisDao.GetPushLogByMeta(ctx, meta)
}

func GetTaskEntryListByMeta(ctx context.Context, meta *model.LogMeta) ([]string, error) {
	return dao.LogRedisDao.GetTaskEntryListByMeta(ctx, meta)
}

func GetTaskStatisticsByID(ctx context.Context, id int) (string, error) {
	if nRecv, err := dao.LogRedisDao.GetTaskStatistics(ctx, id, "receive"); err != nil {
		return "", err
	} else if nShow, err := dao.LogRedisDao.GetTaskStatistics(ctx, id, "show"); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%d received, %d showed", nRecv, nShow), nil
	}
}

func IncrTaskCounter(ctx context.Context, l *model.LogEntry) error {
	return dao.LogRedisDao.IncrTaskCounter(ctx, l.Meta.TaskId, l.Where)
}
