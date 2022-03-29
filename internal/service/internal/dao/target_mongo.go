package dao

import (
	"context"

	"github.com/L-LYR/pns/internal/local_storage"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type _TargetMongoDao struct {
	*internal.TargetMongoDao
}

var (
	TargetMongoDao *_TargetMongoDao
)

func MustInitialize(ctx context.Context) {
	if dao, err := internal.NewTargetMongoDao(ctx); err != nil {
		panic(err)
	} else {
		TargetMongoDao = &_TargetMongoDao{dao}
	}
}

func MustShutdown(ctx context.Context) {
	if err := TargetMongoDao.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func (dao *_TargetMongoDao) SetTarget(
	ctx context.Context,
	t *model.Target,
	opts ...*options.UpdateOptions,
) error {
	collectionName := local_storage.GetAppNameByAppId(t.App.ID)

	filter := bson.D{bson.E{Key: "_id", Value: t.Device.ID}}
	update := bson.D{bson.E{Key: "$set", Value: t}}

	_, err := dao.Collection(collectionName).UpdateOne(ctx, filter, update, opts...)
	return err
}

func (dao *_TargetMongoDao) GetTarget(
	ctx context.Context,
	deviceId string,
	appId int,
) (*model.Target, error) {
	collectionName := local_storage.GetAppNameByAppId(appId)

	result := &model.Target{}
	filter := bson.D{bson.E{Key: "_id", Value: deviceId}}

	if err := dao.
		Collection(collectionName).
		FindOne(ctx, filter).
		Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}
