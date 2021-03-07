package org

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type idKind uint8

const (
	_ idKind = iota
	IDKind_areaID
	IDKind_managerID
	IDKind_metroID
	IDKind_okvedID
	IDKind_locationID
	IDKind_excludeOrgID
)

type ID struct {
	Val  primitive.ObjectID
	Kind idKind
}

func (m Model) GetByIDs(
	ctx context.Context,
	ids []ID,
	skip,
	limit uint32,
) (
	res []Org,
	err error,
) {
	query := bson.M{}
	for _, id := range ids {
		switch id.Kind {
		case IDKind_areaID:
			query["a"] = id.Val
		case IDKind_managerID:
			query["mi"] = id.Val
		case IDKind_metroID:
			query["m.id"] = id.Val
		case IDKind_okvedID:
			query["o"] = id.Val
		case IDKind_locationID:
			query["l"] = id.Val
		case IDKind_excludeOrgID:
			query["_id"] = bson.M{
				"$ne": id.Val,
			}
		}
	}

	cur, err := m.coll.Find(ctx, query, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{
			"_id": -1,
		}))
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}
