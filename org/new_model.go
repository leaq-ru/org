package org

import (
	"github.com/nnqq/scr-org/mongo"
	md "go.mongodb.org/mongo-driver/mongo"
)

func NewModel(db *md.Database) Model {
	return Model{
		coll: db.Collection(mongo.CollOrg),
	}
}
