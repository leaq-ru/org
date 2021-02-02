package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func createIndex(db *mongo.Database) (err error) {
	ctx := context.Background()

	_, err = db.Collection(CollOrg).Indexes().CreateMany(ctx, []mongo.IndexModel{})
	if err != nil {
		return
	}

	return
}
