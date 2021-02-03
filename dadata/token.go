package dadata

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type token struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Value     string             `bson:"v,omitempty"`
	UsedTimes uint32             `bson:"u,omitempty"`
	CreatedAt time.Time          `bson:"ca,omitempty"`
}
