package tools

import (
	"context"
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Initialize(ctx context.Context, mongoNamespaceColl *mongo.Collection) error {

	// Info: Not guarded default namespace for development and testing usage only
	var namespace models.Namespace
	filter := bson.M{
		"name": "default",
	}
	defaultNsResult := mongoNamespaceColl.FindOne(ctx, filter)
	if defaultNsResult.Err() == mongo.ErrNoDocuments {
		namespace = models.Namespace{
			Name:      "default",
			Workflows: make([]string, 0),
			Users:     make([]string, 0),
		}
		res, err := mongoNamespaceColl.InsertOne(ctx, namespace)
		if err != nil {
			return err
		}
		fmt.Printf("Inserted default namespace: %v\n", res.InsertedID)
	} else {
		err := defaultNsResult.Decode(&namespace)
		if err != nil {
			return err
		}

		fmt.Printf("Default namespace already exists %s\n", namespace.Id)
	}

	// Loading installable integrations
	_, err := integrations.GetInstallableIntegrationsLib()
	if err != nil {
		return fmt.Errorf("cannot parse installable error: %s", err)
	}

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
