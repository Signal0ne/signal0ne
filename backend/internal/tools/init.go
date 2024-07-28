package tools

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Initialize(ctx context.Context, mongoNamespaceColl *mongo.Collection) error {
	namespace := bson.M{
		"name":      "default",
		"workflows": make([]string, 0),
		"users":     make([]string, 0),
	}
	res, err := mongoNamespaceColl.InsertOne(ctx, namespace)
	if err != nil {
		return err
	}

	fmt.Printf("Inserted default namespace: %v\n", res.InsertedID)

	return nil
}

func InitMongoClient(ctx context.Context, mongoUri string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
