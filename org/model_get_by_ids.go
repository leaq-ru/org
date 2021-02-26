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
		var key string
		switch id.Kind {
		case IDKind_areaID:
			key = "a"
		case IDKind_managerID:
			key = "mi"
		case IDKind_metroID:
			key = "m.id"
		case IDKind_okvedID:
			key = "o"
		case IDKind_locationID:
			key = "l"
		}
		query[key] = id.Val
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
