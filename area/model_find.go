package area

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
	name, fiasID, kladrID, typ, typeFull string,
	kind Kind,
) (
	id primitive.ObjectID,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	upsertID := primitive.NewObjectID()
	s := slug.Make(name)

	ur, err := m.coll.UpdateOne(ctx, area{
		Slug: s,
	}, bson.M{
		"$setOnInsert": area{
			ID:       upsertID,
			FiasID:   fiasID,
			KladrID:  kladrID,
			Name:     name,
			Kind:     kind,
			Type:     typ,
			TypeFull: typeFull,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return
	}

	if ur != nil && ur.UpsertedCount == 1 {
		id = upsertID
		return
	}

	var doc area
	err = m.coll.FindOne(ctx, area{
		Slug: s,
	}).Decode(&doc)
	if err != nil {
		return
	}

	id = doc.ID
	return
}
