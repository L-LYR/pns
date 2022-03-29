package internal

import (
	"context"

	"github.com/L-LYR/pns/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_MongoDaoConfigName = "mongo"
)

type TargetMongoDao struct {
	dbName string
	client mongo.Client
}

func NewTargetMongoDao(ctx context.Context) (*TargetMongoDao, error) {
	cfg := config.MustLoadMongoDaoConfig(ctx, _MongoDaoConfigName)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI()))
	if err != nil {
		return nil, err
	}
	return &TargetMongoDao{
		client: *client,
		dbName: cfg.Name,
	}, nil
}

func (dao *TargetMongoDao) DB() *mongo.Database {
	return dao.client.Database(dao.dbName)
}

func (dao *TargetMongoDao) Collection(name string) *mongo.Collection {
	return dao.DB().Collection(name)
}

func (dao *TargetMongoDao) Shutdown(ctx context.Context) error {
	return dao.client.Disconnect(ctx)
}
