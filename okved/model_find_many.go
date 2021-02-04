package okved

import (
	"context"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"time"
)

type Upsert struct {
	Code         string
	Name         string
	CodeWithName string
	Kind         kind
}

func (m Model) FindMany(
	ctx context.Context,
	vals []Upsert,
) (
	ids []primitive.ObjectID,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ids = make([]primitive.ObjectID, len(vals))

	var eg errgroup.Group
	for _i, _val := range vals {
		i := _i
		val := _val
		eg.Go(func() error {
			upsertID := primitive.NewObjectID()
			s := slug.Make(val.CodeWithName)

			ur, e := m.coll.UpdateOne(ctx, okved{
				Slug: s,
			}, bson.M{
				"$setOnInsert": okved{
					ID:           upsertID,
					Code:         val.Code,
					Name:         val.Name,
					CodeWithName: val.CodeWithName,
					Kind:         val.Kind,
				},
			}, options.Update().SetUpsert(true))
			if e != nil {
				return e
			}

			if ur != nil && ur.UpsertedCount == 1 {
				ids[i] = upsertID
				return nil
			}

			var doc okved
			e = m.coll.FindOne(ctx, okved{
				Slug: s,
			}).Decode(&doc)
			if e != nil {
				return e
			}

			ids[i] = doc.ID
			return nil
		})
	}
	err = eg.Wait()
	return
}
