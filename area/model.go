package area

import "go.mongodb.org/mongo-driver/mongo"

type Model struct {
	coll *mongo.Collection
}
