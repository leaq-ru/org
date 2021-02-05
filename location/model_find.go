package location

import (
	"context"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (m Model) Find(
	ctx context.Context,
	name string,
) (
	id primitive.ObjectID,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if name == "" {
		return
	}

	upsertID := primitive.NewObjectID()
	s := slug.Make(name)

	ur, err := m.coll.UpdateOne(ctx, location{
		Slug: s,
	}, bson.M{
		"$setOnInsert": location{
			ID:   upsertID,
			Name: name,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return
	}

	if ur != nil && ur.UpsertedCount == 1 {
		id = upsertID
		return
	}

	var doc location
	err = m.coll.FindOne(ctx, location{
		Slug: s,
	}).Decode(&doc)
	if err != nil {
		return
	}

	id = doc.ID
	return
}
