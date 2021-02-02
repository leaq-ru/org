package metro

import "go.mongodb.org/mongo-driver/bson/primitive"

type metro struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Slug   string             `bson:"s,omitempty"`
	AreaID primitive.ObjectID `bson:"a,omitempty"`
	Name   string             `bson:"n,omitempty"`
	Line   string             `bson:"l,omitempty"`
}
