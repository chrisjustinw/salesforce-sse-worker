package model

import "go.mongodb.org/mongo-driver/v2/bson"

type ConversationMapping struct {
	Id        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Token     string        `json:"token" bson:"token"`
	Partition int           `json:"partition" bson:"partition"`
}
