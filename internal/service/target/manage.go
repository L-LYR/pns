package target

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/cache"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"go.mongodb.org/mongo-driver/mongo"
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
	appId int,
	deviceId string,
) (*model.Target, error) {
	appName, ok := cache.Config.GetAppNameByAppId(appId)
	if !ok {
		return nil, errors.New("wrong app id")
	}
	return dao.TargetMongoDao.GetTarget(ctx, appName, deviceId)
}

func UpdateToken(
	ctx context.Context,
	appId int,
	deviceId string,
	tokenName string,
	token string,
) error {
	appName, ok := cache.Config.GetAppNameByAppId(appId)
	if !ok {
		return errors.New("wrong app id")
	}
	return dao.TargetMongoDao.SetTargetToken(
		ctx,
		appName,
		deviceId,
		tokenName,
		token,
	)
}

func Cursor(ctx context.Context, appName string) (*mongo.Cursor, error) {
	return dao.TargetMongoDao.NaiveCursor(ctx, appName)
}
