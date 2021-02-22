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
	IDKind_AreaID
	IDKind_ManagerID
	IDKind_MetroID
	IDKind_OkvedID
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
		case IDKind_AreaID:
			key = "a"
		case IDKind_ManagerID:
			key = "mi"
		case IDKind_MetroID:
			key = "m.id"
		case IDKind_OkvedID:
			key = "o"
		}
		query[key] = id.Val
	}

	cur, err := m.coll.Find(ctx, query, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}
