package dadata

import (
	"github.com/nnqq/scr-org/mongo"
	md "go.mongodb.org/mongo-driver/mongo"
)

func newTokenPool(tokens []string, db *md.Database) tokenPool {
	return tokenPool{
		tokens: tokens,
		coll:   db.Collection(mongo.CollDaDataToken),
	}
}
