package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoConfigName = "database.pns_mongo"
)

type TargetMongoDao struct {
	dbName string
	client mongo.Client
}

func NewTargetMongoDao(ctx context.Context) (*TargetMongoDao, error) {
	var cfg map[string]interface{}
	if g.Cfg().Available(ctx) {
		cfg = g.Cfg().MustGet(ctx, MongoConfigName).Map()
	}
	uriFields := make([]interface{}, 0, 4)
	for _, field := range []string{"user", "pass", "host", "port"} {
		if _, ok := cfg[field]; !ok {
			return nil, fmt.Errorf("%s is invalid in mongo config", field)
		}
		uriFields = append(uriFields, cfg[field])
	}
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", uriFields...)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	dao := &TargetMongoDao{client: *client}
	if dbName, ok := cfg["name"].(string); ok {
		dao.dbName = dbName
	} else {
		return nil, errors.New("db name is unknown")
	}
	return dao, nil
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
