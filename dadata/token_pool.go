package dadata

import "go.mongodb.org/mongo-driver/mongo"

type TokenPool struct {
	tokens []string
	coll   *mongo.Collection
}
