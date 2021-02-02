package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strings"
	"time"
)

func NewConn(ctx context.Context, serviceName string, url string) (db *mongo.Database, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := mongo.Connect(ctx, options.Client().
		SetWriteConcern(writeconcern.New(
			writeconcern.WMajority(),
			writeconcern.J(true),
		)).
		SetReadConcern(readconcern.Majority()).
		SetReadPreference(readpref.SecondaryPreferred()).
		ApplyURI(url))
	if err != nil {
		return
	}

	err = conn.Ping(ctx, nil)
	if err != nil {
		return
	}

	db = conn.Database(strings.ReplaceAll(serviceName, "-", "_"))

	err = createIndex(db)
	return
}
