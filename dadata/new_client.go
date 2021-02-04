package dadata

import "go.mongodb.org/mongo-driver/mongo"

func NewClient(tokens []string, db *mongo.Database) Client {
	return Client{
		tp: newTokenPool(tokens, db),
	}
}
