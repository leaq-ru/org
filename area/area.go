package area

import "go.mongodb.org/mongo-driver/bson/primitive"

type Area struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Slug     string             `bson:"s,omitempty"`
	FiasID   string             `bson:"fi,omitempty"`
	KladrID  string             `bson:"ki,omitempty"`
	Name     string             `bson:"n,omitempty"`
	Kind     Kind               `bson:"k,omitempty"`
	Type     string             `bson:"t,omitempty"`
	TypeFull string             `bson:"tf,omitempty"`
}

type Kind uint8

const (
	_ Kind = iota
	Kind_city
	Kind_settlement
)
