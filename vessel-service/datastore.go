package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	dbName = "shipper"
	vesselCollection = "vessels"
)

func CreateClient(uri string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}

func CreateConnection(uri string, dbName string, collection string) (*mongo.Collection, error) {
	client, err := CreateClient(uri)
	if err != nil {
		return nil, err
	}
	return client.Database(dbName).Collection(collection), nil
}
