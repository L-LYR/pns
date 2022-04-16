package target

import (
	"context"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/cache"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Upsert(ctx context.Context, target *model.Target) error {
	now := time.Now()
	appName, _ := cache.Config.GetAppNameByAppId(target.App.ID)
	if result, _ := dao.TargetMongoDao.GetTarget( // ignore this error
		ctx,
		appName,
		target.Device.ID,
	); result != nil {
		if result.Equal(target) {
			return nil
		}
		target.Tokens = result.Tokens
		target.TokenUpdateTime = result.TokenUpdateTime
		target.CreateTime = result.CreateTime
	} else {
		target.CreateTime = now
	}

	target.LastActiveTime = now
	target.InfoUpdateTime = now

	if err := dao.TargetMongoDao.SetTarget(
		ctx,
		appName,
		target,
		options.Update().SetUpsert(true),
	); err != nil {
		return err
	}
	return nil
}

func Query(
	ctx context.Context,
	appName string,
	deviceId string,
) (*model.Target, error) {
	return dao.TargetMongoDao.GetTarget(ctx, appName, deviceId)
}

func UpdateToken(
	ctx context.Context,
	appName string,
	deviceId string,
	tokenName string,
	token string,
) error {
	return dao.TargetMongoDao.SetTargetToken(
		ctx,
		appName,
		deviceId,
		tokenName,
		token,
	)
}

func Scan(
	ctx context.Context,
	appName string,
	fn func(*model.Target) error,
	errorHandler func(error),
) (int64, error) {
	cursor, err := dao.TargetMongoDao.NaiveCursor(ctx, appName)
	if err != nil {
		return 0, err
	}
	count := int64(0)
	for cursor.Next(context.TODO()) {
		result := &model.Target{}
		if err := cursor.Decode(&result); err != nil {
			errorHandler(err)
			continue
		}
		if err := fn(result); err != nil {
			errorHandler(err)
			continue
		}
		count++
	}
	return count, cursor.Close(ctx)
}
