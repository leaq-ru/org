package okved

import "go.mongodb.org/mongo-driver/bson/primitive"

type okved struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Slug         string             `bson:"s,omitempty"`
	Code         string             `bson:"c,omitempty"`
	Name         string             `bson:"n,omitempty"`
	CodeWithName string             `bson:"cn,omitempty"`
	Kind         kind               `bson:"k,omitempty"`
}

type kind uint8

const (
	_ kind = iota
	Kind_y2001
	Kind_y2004
)
