package area

import "go.mongodb.org/mongo-driver/bson/primitive"

type area struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Slug     string             `bson:"s,omitempty"`
	FiasID   string             `bson:"fi,omitempty"`
	KladrID  string             `bson:"ki,omitempty"`
	Name     string             `bson:"n,omitempty"`
	Kind     kind               `bson:"k,omitempty"`
	Type     string             `bson:"t,omitempty"`
	TypeFull string             `bson:"tf,omitempty"`
}

type kind uint8

const (
	_ kind = iota
	kind_city
	kind_settlement
)
