package manager

import "go.mongodb.org/mongo-driver/bson/primitive"

type Manager struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Slug string             `bson:"s,omitempty"`
	Name string             `bson:"n,omitempty"`
}
