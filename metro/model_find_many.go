package metro

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

type FindManyReqItem struct {
	Name     string
	Line     string
	Distance float32
}

type FindManyResItem struct {
	ID       primitive.ObjectID
	Distance float32
}

func (m Model) FindMany(
	ctx context.Context,
	areaID primitive.ObjectID,
	vals []FindManyReqItem,
) (
	res []FindManyResItem,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res = make([]FindManyResItem, len(vals))
	var eg errgroup.Group
	for _i, _val := range vals {
		i := _i
		val := _val
		eg.Go(func() error {
			upsertID := primitive.NewObjectID()

			forSlug := val.Name
			if val.Line != "" {
				forSlug = strings.Join([]string{
					val.Line,
					val.Name,
				}, " ")
			}
			s := slug.Make(forSlug)

			ur, e := m.coll.UpdateOne(ctx, metro{
				Slug: s,
			}, bson.M{
				"$setOnInsert": metro{
					ID:     upsertID,
					AreaID: areaID,
					Name:   val.Name,
					Line:   val.Line,
				},
			}, options.Update().SetUpsert(true))
			if e != nil {
				return e
			}

			if ur != nil && ur.UpsertedCount == 1 {
				res[i] = FindManyResItem{
					ID:       upsertID,
					Distance: val.Distance,
				}
				return nil
			}

			var doc metro
			e = m.coll.FindOne(ctx, metro{
				Slug: s,
			}).Decode(&doc)
			if e != nil {
				return e
			}

			res[i] = FindManyResItem{
				ID:       doc.ID,
				Distance: val.Distance,
			}
			return nil
		})
	}
	err = eg.Wait()
	return
}
