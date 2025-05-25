package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	BaseRepository interface {
		FindOne(ctx context.Context, collection string, findQuery map[string]interface{}) *mongo.SingleResult
		ReplaceOne(ctx context.Context, collection string, query interface{}, data interface{}) (result *mongo.UpdateResult, err error)
	}

	BaseRepositoryImpl struct {
		MongoDatabase *mongo.Database
	}
)

func NewBaseRepository(mongoDatabase *mongo.Database) BaseRepository {
	return &BaseRepositoryImpl{MongoDatabase: mongoDatabase}
}

func (m *BaseRepositoryImpl) FindOne(ctx context.Context, collection string, query map[string]interface{}) *mongo.SingleResult {
	return m.MongoDatabase.Collection(collection).FindOne(ctx, query)
}

func (m *BaseRepositoryImpl) ReplaceOne(ctx context.Context, collection string, query interface{}, data interface{}) (result *mongo.UpdateResult, err error) {
	return m.MongoDatabase.Collection(collection).ReplaceOne(ctx, query, data, options.Replace().SetUpsert(true))
}
