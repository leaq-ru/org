package okved

import (
	"context"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

type Upsert struct {
	Code string
	Name string
	Kind string
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

	if len(vals) == 0 {
		return
	}

	ids = make([]primitive.ObjectID, len(vals))

	var eg errgroup.Group
	for _i, _val := range vals {
		i := _i
		val := _val
		eg.Go(func() error {
			upsertID := primitive.NewObjectID()
			codeWithName := strings.Join([]string{
				val.Code,
				val.Name,
			}, " ")
			s := slug.Make(codeWithName)

			ur, e := m.coll.UpdateOne(ctx, okved{
				Slug: s,
			}, bson.M{
				"$setOnInsert": okved{
					ID:           upsertID,
					Code:         val.Code,
					Name:         val.Name,
					CodeWithName: codeWithName,
					Kind:         toKind(val.Kind),
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

func toKind(in string) kind {
	switch in {
	case "2001":
		return kind_y2001
	case "2014":
		return kind_y2014
	default:
		return 0
	}
}
