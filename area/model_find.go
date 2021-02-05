package area

import (
	"context"
	"fmt"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
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

	upsertOK := false
	s := slug.Make(name)
	for i := 1; i <= 100; i += 1 {
		upsertID := primitive.NewObjectID()
		ur, e := m.coll.UpdateOne(ctx, area{
			FiasID: fiasID,
		}, bson.M{
			"$setOnInsert": area{
				ID:       upsertID,
				Slug:     s,
				KladrID:  kladrID,
				Name:     name,
				Kind:     kind,
				Type:     typ,
				TypeFull: typeFull,
			},
		}, options.Update().SetUpsert(true))
		if e != nil {
			s += "-" + strconv.Itoa(i)
			continue
		}

		if ur != nil && ur.UpsertedCount == 1 {
			id = upsertID
			return
		}
		upsertOK = true
		break
	}

	if !upsertOK {
		err = fmt.Errorf("failed to upsert area. fiasID=%s, name=%s", fiasID, name)
		return
	}

	var doc area
	err = m.coll.FindOne(ctx, area{
		FiasID: fiasID,
	}).Decode(&doc)
	if err != nil {
		return
	}

	id = doc.ID
	return
}
