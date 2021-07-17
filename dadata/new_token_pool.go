package dadata

import (
	"github.com/leaq-ru/org/mongo"
	md "go.mongodb.org/mongo-driver/mongo"
)

func newTokenPool(tokens []string, db *md.Database) tokenPool {
	return tokenPool{
		tokens: tokens,
		coll:   db.Collection(mongo.CollDaDataToken),
	}
}
