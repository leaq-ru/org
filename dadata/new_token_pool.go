package dadata

import (
	"github.com/nnqq/scr-org/mongo"
	md "go.mongodb.org/mongo-driver/mongo"
)

func NewTokenPool(tokens []string, db *md.Database) TokenPool {
	return TokenPool{
		tokens: tokens,
		coll:   db.Collection(mongo.CollDaDataToken),
	}
}
