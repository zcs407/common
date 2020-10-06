package conn

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	conn *mongo.Database
}

func NewMongoClient(url, dbName string) (m *MongoClient, err error) {
	ctx := context.Background()
	m = &MongoClient{}
	conn, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return
	}

	err = conn.Connect(ctx)
	if err != nil {
		return
	}
	m.conn = conn.Database(dbName)
	return
}
