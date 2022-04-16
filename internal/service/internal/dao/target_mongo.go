package dao

import (
	"context"
	"time"

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
	TargetMongoDao _TargetMongoDao
)

func MustInitTargetMongoDao(ctx context.Context) {
	if dao, err := internal.NewTargetMongoDao(ctx); err != nil {
		panic(err)
	} else {
		TargetMongoDao = _TargetMongoDao{dao}
	}
}

func MustShutdownTargetMongoDao(ctx context.Context) {
	if err := TargetMongoDao.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func (dao *_TargetMongoDao) SetTarget(
	ctx context.Context,
	appName string,
	t *model.Target,
	opts ...*options.UpdateOptions,
) error {
	filter := bson.D{bson.E{Key: "_id", Value: t.Device.ID}}
	update := bson.D{bson.E{Key: "$set", Value: t}}

	_, err := dao.Collection(appName).UpdateOne(ctx, filter, update, opts...)

	return err
}

func (dao *_TargetMongoDao) SetTargetToken(
	ctx context.Context,
	appName string,
	deviceId string,
	tokenName string,
	token string,
) error {
	now := time.Now()
	filter := bson.D{bson.E{Key: "_id", Value: deviceId}}
	update := bson.D{bson.E{
		Key: "$set",
		Value: bson.D{
			bson.E{Key: "tokens", Value: map[string]string{tokenName: token}},
			bson.E{Key: "lastActiveTime", Value: now},
			bson.E{Key: "tokenUpdateTime", Value: now},
		},
	}}
	_, err := dao.Collection(appName).UpdateOne(ctx, filter, update)
	return err
}

func (dao *_TargetMongoDao) GetTarget(
	ctx context.Context,
	appName string,
	deviceId string,
) (*model.Target, error) {
	result := &model.Target{}
	filter := bson.D{bson.E{Key: "_id", Value: deviceId}}

	if err := dao.
		Collection(appName).
		FindOne(ctx, filter).
		Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (dao *_TargetMongoDao) CountTarget(ctx context.Context, appName string) (int64, error) {
	return dao.TargetMongoDao.Collection(appName).CountDocuments(ctx, bson.D{})
}

func (dao *_TargetMongoDao) NaiveCursor(ctx context.Context, appName string) (*mongo.Cursor, error) {
	return dao.TargetMongoDao.Collection(appName).Find(ctx, bson.D{})
}
