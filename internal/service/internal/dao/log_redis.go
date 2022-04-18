package dao

import (
	"context"
	"strconv"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao/internal"
	"github.com/go-redis/redis/v8"
)

type _LogRedisDao struct {
	*internal.LogRedisDao
}

var (
	LogRedisDao _LogRedisDao
)

func MustInitLogRedisDao(ctx context.Context) {
	if dao, err := internal.NewLogRedisDao(ctx); err != nil {
		panic(err)
	} else {
		LogRedisDao = _LogRedisDao{dao}
	}
}

func MustShutdownLogRedisDao(ctx context.Context) {
	if err := LogRedisDao.Shutdown(ctx); err != nil {
		panic(err)
	}
}

/*

1. Task Log

Task ID -> Request Log / Send Log / Recv Log / Show Log ......

2. Task Entry List

App ID + Device ID -> Task ID 1 / Task ID 2 / Task ID 3 ......

*/

// TODO: make this configurable
const (
	_InactiveDuration = 3 * 24 * 3600 * time.Second
)

func (dao *_LogRedisDao) _SweepExpireEntry(
	ctx context.Context,
	key string,
) error {
	upperBound := time.Now().Add(-_InactiveDuration).Unix()
	if _, err := dao.LogRedisDao.Client(ctx).ZRemRangeByScore(
		ctx,
		key,
		"-inf",
		strconv.FormatInt(upperBound, 10),
	).Result(); err != nil {
		return err
	}
	return nil
}

func (dao *_LogRedisDao) _ZListAppend(
	ctx context.Context,
	key string,
	value string,
	score float64,
) error {
	client := dao.LogRedisDao.Client(ctx)
	if _, err := client.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: value,
	}).Result(); err != nil {
		return err
	}
	if _, err := client.
		Expire(ctx, key, _InactiveDuration).
		Result(); err != nil {
		return err
	}
	return nil
}

func (dao *_LogRedisDao) AppendTaskLog(
	ctx context.Context,
	log *model.LogEntry,
) error {
	value, err := log.Encode()
	if err != nil {
		return err
	}
	key := log.Meta.TaskStatusKey()
	if err := dao._ZListAppend(ctx, key, value, float64(log.T)); err != nil {
		return err
	}
	if err := dao._SweepExpireEntry(ctx, key); err != nil {
		return err
	}
	return nil
}

func (dao *_LogRedisDao) AppendPushLog(
	ctx context.Context,
	log *model.LogEntry,
) error {
	value, err := log.Encode()
	if err != nil {
		return err
	}
	key := log.Meta.PushKey()
	if err := dao._ZListAppend(ctx, key, value, float64(log.T)); err != nil {
		return err
	}
	if err := dao._SweepExpireEntry(ctx, key); err != nil {
		return err
	}
	return nil
}

func (dao *_LogRedisDao) AppendTaskEntry(
	ctx context.Context,
	meta *model.LogMeta,
) error {
	key := meta.EntryKey()
	value := meta.TaskStatusKey()
	if err := dao._ZListAppend(ctx, key, value, float64(time.Now().UnixMilli())); err != nil {
		return err
	}
	if err := dao._SweepExpireEntry(ctx, key); err != nil {
		return err
	}
	return nil
}

func (dao *_LogRedisDao) GetTaskLogByID(
	ctx context.Context,
	id int,
) ([]*model.LogEntry, error) {
	client := dao.LogRedisDao.Client(ctx)
	key := strconv.FormatInt(int64(id), 10)
	rawLog, err := client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: "inf",
	}).Result()
	if err != nil {
		return nil, err
	}
	log := make([]*model.LogEntry, 0, len(rawLog))
	for i := range rawLog {
		entry := &model.LogEntry{}
		if err := entry.Decode(rawLog[i]); err != nil {
			log = append(log, model.DummyEntry)
		} else {
			log = append(log, entry)
		}
	}
	return log, nil
}

func (dao *_LogRedisDao) GetTaskLastLogByID(
	ctx context.Context,
	id int,
) (*model.LogEntry, error) {
	client := dao.LogRedisDao.Client(ctx)
	key := strconv.FormatInt(int64(id), 10)
	rawLog, err := client.ZRevRange(ctx, key, 0, 0).Result()
	if err != nil {
		return nil, err
	}
	if len(rawLog) == 0 {
		return nil, nil
	}
	entry := &model.LogEntry{}
	if err := entry.Decode(rawLog[0]); err != nil {
		return model.DummyEntry, err
	}
	return entry, nil
}

func (dao *_LogRedisDao) GetPushLogByMeta(
	ctx context.Context,
	meta *model.LogMeta,
) ([]*model.LogEntry, error) {
	client := dao.LogRedisDao.Client(ctx)
	key := meta.PushKey()
	rawLog, err := client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: "inf",
	}).Result()
	if err != nil {
		return nil, err
	}
	log := make([]*model.LogEntry, 0, len(rawLog))
	for i := range rawLog {
		entry := &model.LogEntry{}
		if err := entry.Decode(rawLog[i]); err != nil {
			log = append(log, model.DummyEntry)
		} else {
			log = append(log, entry)
		}
	}
	return log, nil
}
