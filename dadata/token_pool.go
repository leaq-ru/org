package dadata

import "go.mongodb.org/mongo-driver/mongo"

type tokenPool struct {
	tokens []string
	coll   *mongo.Collection
}
