package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func createIndex(db *mongo.Database) (err error) {
	ctx := context.Background()

	_, err = db.Collection(CollOrg).Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"s": 1,
		},
		Options: options.Index().SetUnique(true),
	},
	//	{Keys: bson.D{{
	//		Key:   "a",
	//		Value: 1,
	//	}, {
	//		Key:   "_id",
	//		Value: -1,
	//	}},
	//}, {
	//	Keys: bson.D{{
	//		Key:   "l",
	//		Value: 1,
	//	}, {
	//		Key:   "_id",
	//		Value: -1,
	//	}},
	//}, {
	//	Keys: bson.D{{
	//		Key:   "mi",
	//		Value: 1,
	//	}, {
	//		Key:   "_id",
	//		Value: -1,
	//	}},
	//}, {
	//	Keys: bson.D{{
	//		Key:   "o",
	//		Value: 1,
	//	}, {
	//		Key:   "_id",
	//		Value: -1,
	//	}},
	//}, {
	//	Keys: bson.D{{
	//		Key:   "m.id",
	//		Value: 1,
	//	}, {
	//		Key:   "_id",
	//		Value: -1,
	//	}},
	//}, {
	//	Keys: bson.D{{
	//		Key:   "sk",
	//		Value: 1,
	//	}, {
	//		Key:   "_id",
	//		Value: -1,
	//	}},
	//}, {
	//	Keys: bson.D{{
	//		Key:   "rd",
	//		Value: -1,
	//	}, {
	//		Key:   "_id",
	//		Value: -1,
	//	}},
	//}
	})
	if err != nil {
		return
	}

	_, err = db.Collection(CollArea).Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"s": 1,
		},
		Options: options.Index().SetUnique(true),
	}, {
		Keys: bson.M{
			"fi": 1,
		},
		Options: options.Index().SetUnique(true),
	}})
	if err != nil {
		return
	}

	_, err = db.Collection(CollLocation).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"s": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return
	}

	_, err = db.Collection(CollManager).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"s": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return
	}

	_, err = db.Collection(CollOkved).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"s": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return
	}

	_, err = db.Collection(CollMetro).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"s": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return
	}
	_, err = db.Collection(CollDaDataToken).Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"v": 1,
		},
		Options: options.Index().SetUnique(true),
	}, {
		Keys: bson.M{
			"u": 1,
		},
	}, {
		Keys: bson.M{
			"ca": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(int32((24 * time.Hour).Seconds())),
	}})

	return
}
