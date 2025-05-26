package library

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"salesforce-sse-worker/configs"
)

type (
	MongoDatabase interface {
		Find(ctx context.Context, collection string, findQuery map[string]interface{}) (*mongo.Cursor, error)
		FindOne(ctx context.Context, collection string, findQuery map[string]interface{}) *mongo.SingleResult
		ReplaceOne(ctx context.Context, collection string, query interface{}, data interface{}) (result *mongo.UpdateResult, err error)
	}

	MongoDatabaseImpl struct {
		db *mongo.Database
	}
)

func NewMongoDatabase(cfg configs.MongoConfig, client *mongo.Client) MongoDatabase {
	return &MongoDatabaseImpl{db: client.Database(cfg.DatabaseName)}
}

func (m *MongoDatabaseImpl) Find(ctx context.Context, collection string, query map[string]interface{}) (*mongo.Cursor, error) {
	return m.db.Collection(collection).Find(ctx, query)
}

func (m *MongoDatabaseImpl) FindOne(ctx context.Context, collection string, query map[string]interface{}) *mongo.SingleResult {
	return m.db.Collection(collection).FindOne(ctx, query)
}

func (m *MongoDatabaseImpl) ReplaceOne(ctx context.Context, collection string, query interface{}, data interface{}) (result *mongo.UpdateResult, err error) {
	return m.db.Collection(collection).ReplaceOne(ctx, query, data, options.Replace().SetUpsert(true))
}
