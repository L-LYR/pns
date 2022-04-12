package target

import (
	"context"
	"time"

	"github.com/L-LYR/pns/internal/local_storage"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Upsert(ctx context.Context, target *model.Target) error {
	now := time.Now()
	appName, _ := local_storage.GetAppNameByAppId(target.App.ID)
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
