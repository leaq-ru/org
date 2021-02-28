package area

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m Model) GetByID(ctx context.Context, id primitive.ObjectID) (res Area, err error) {
	err = m.coll.FindOne(ctx, Area{
		ID: id,
	}).Decode(&res)
	return
}
