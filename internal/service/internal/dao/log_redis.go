package dao

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/L-LYR/pns/internal/config"
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

Task ID -> Server Log List

2. Push Log

App ID : Device ID : Task ID -> Client Log List

3. Task Entry List

1) App ID : Device ID -> Task ID 1 / Task ID 2 / Task ID 3 ......

2) App ID -> Broadcast Task ID 1 / ......
*/

func (dao *_LogRedisDao) _AppendLog(
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
		Expire(ctx, key, config.LogExpireTime()).
		Result(); err != nil {
		return err
	}
	upperBound := time.Now().Add(-config.LogExpireTime()).Unix()
	if _, err := client.ZRemRangeByScore(
		ctx, key, "-inf",
		strconv.FormatInt(upperBound, 10),
	).Result(); err != nil {
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
	return dao._AppendLog(ctx, log.Meta.TaskStatusKey(), value, float64(log.T))
}

func (dao *_LogRedisDao) AppendPushLog(
	ctx context.Context,
	log *model.LogEntry,
) error {
	value, err := log.Encode()
	if err != nil {
		return err
	}
	return dao._AppendLog(ctx, log.Meta.PushKey(), value, float64(log.T))
}

func (dao *_LogRedisDao) AppendTaskEntry(
	ctx context.Context,
	meta *model.LogMeta,
) error {
	return dao._AppendLog(
		ctx, meta.EntryKey(),
		meta.TaskStatusKey(),
		float64(time.Now().UnixMilli()),
	)
}

func (dao *_LogRedisDao) CheckAndAppendTaskEntry(
	ctx context.Context,
	meta *model.LogMeta,
) error {
	_, err := dao.Client(ctx).ZScore(
		ctx, meta.EntryKey(),
		meta.TaskStatusKey(),
	).Result()
	if errors.Is(err, redis.Nil) {
		return dao.AppendTaskEntry(ctx, meta)
	}
	return nil
}

func (dao *_LogRedisDao) GetTaskLogByID(
	ctx context.Context,
	id int64,
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
	id int64,
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
	rawLog, err := client.ZRangeByScore(ctx, meta.PushKey(), &redis.ZRangeBy{
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

func (dao *_LogRedisDao) GetTaskEntryListByMeta(
	ctx context.Context,
	meta *model.LogMeta,
) ([]string, error) {
	client := dao.LogRedisDao.Client(ctx)
	rawLog, err := client.ZRangeByScore(ctx, meta.EntryKey(), &redis.ZRangeBy{
		Min: "-inf",
		Max: "inf",
	}).Result()
	if err != nil {
		return nil, err
	}
	return rawLog, nil
}

func (dao *_LogRedisDao) IncrTaskCounter(
	ctx context.Context,
	taskId int64,
	event string, // receive or show
) error {
	client := dao.LogRedisDao.Client(ctx)
	key := fmt.Sprintf("%d:%s", taskId, event)
	if _, err := client.Incr(ctx, key).Result(); err != nil {
		return err
	}
	if _, err := client.
		Expire(ctx, key, config.LogExpireTime()).
		Result(); err != nil {
		return err
	}
	return nil
}

func (dao *_LogRedisDao) GetTaskStatistics(
	ctx context.Context,
	taskId int64,
	event string,
) (int64, error) {
	key := fmt.Sprintf("%d:%s", taskId, event)
	return dao.LogRedisDao.Client(ctx).Get(ctx, key).Int64()
}
