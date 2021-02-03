package dadata

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	md "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"time"
)

const dailyLimitPerToken = 10000

var ErrNoFreeTokens = errors.New("no tokens which not exceeded daily limit")

func (t TokenPool) Get(ctx context.Context) (tok string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var eg errgroup.Group
	for _, _item := range t.tokens {
		item := _item
		eg.Go(func() error {
			_, e := t.coll.UpdateOne(ctx, token{
				Value: item,
			}, bson.M{
				"$setOnInsert": token{
					CreatedAt: time.Now().UTC(),
				},
			}, options.Update().SetUpsert(true))
			return e
		})
	}
	err = eg.Wait()
	if err != nil {
		return
	}

	var doc token
	err = t.coll.FindOneAndUpdate(ctx, bson.M{
		"u": bson.M{
			"$lt": dailyLimitPerToken,
		},
	}, bson.M{
		"$inc": token{
			UsedTimes: 1,
		},
	}).Decode(&doc)
	if err != nil {
		if errors.Is(err, md.ErrNoDocuments) {
			err = ErrNoFreeTokens
		}
		return
	}

	// to not flood dadata.ru
	time.Sleep(75 * time.Millisecond)
	tok = doc.Value
	return
}
