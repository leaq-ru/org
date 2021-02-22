package okved

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m Model) GetByIDs(ctx context.Context, ids []primitive.ObjectID) (res []Okved, err error) {
	cur, err := m.coll.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return
	}

	err = cur.Decode(&res)
	return
}
