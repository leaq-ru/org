package location

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Slug string             `bson:"s,omitempty"`
	Name string             `bson:"n,omitempty"`
}
