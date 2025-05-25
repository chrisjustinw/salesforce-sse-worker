package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"salesforce-sse-worker/internal/model"
)

const (
	conversationMapping = "conversation_mapping"
)

type (
	ConversationMappingRepository interface {
		Upsert(ctx context.Context, data model.ConversationMapping) (*mongo.UpdateResult, error)
		FindAll(ctx context.Context) ([]model.ConversationMapping, error)
		FindOneByPartition(ctx context.Context, partition int) (*model.ConversationMapping, error)
	}

	ConversationMappingRepositoryImpl struct {
		MongoDatabase BaseRepository
	}
)

func NewConversationMappingRepository(mongoDatabase BaseRepository) ConversationMappingRepository {
	return &ConversationMappingRepositoryImpl{
		MongoDatabase: mongoDatabase,
	}
}

func (s *ConversationMappingRepositoryImpl) Upsert(ctx context.Context, data model.ConversationMapping) (*mongo.UpdateResult, error) {
	query := map[string]interface{}{
		"partition": data.Partition,
	}

	return s.MongoDatabase.ReplaceOne(ctx, conversationMapping, query, data)
}

func (s *ConversationMappingRepositoryImpl) FindAll(ctx context.Context) ([]model.ConversationMapping, error) {
	query := map[string]interface{}{}

	cursor, err := s.MongoDatabase.Find(ctx, conversationMapping, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.ConversationMapping
	for cursor.Next(ctx) {
		var mapping model.ConversationMapping
		if err := cursor.Decode(&mapping); err != nil {
			return nil, err
		}
		results = append(results, mapping)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *ConversationMappingRepositoryImpl) FindOneByPartition(ctx context.Context, partition int) (*model.ConversationMapping, error) {
	query := map[string]interface{}{
		"partition": partition,
	}

	singleResult := s.MongoDatabase.FindOne(ctx, conversationMapping, query)
	if err := singleResult.Err(); err != nil {
		return nil, err
	}

	var result model.ConversationMapping
	if err := singleResult.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
