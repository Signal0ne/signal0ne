package tools

import (
	"context"
	"fmt"
	"signal0ne/internal/models" //only internal import allowed

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func OnboardAdmin(
	ctx context.Context,
	mongoNamespaceColl *mongo.Collection,
	admin models.User) error {

	var namespace models.Namespace
	filter := bson.M{
		"name": "default",
	}
	defaultNsResult := mongoNamespaceColl.FindOne(ctx, filter)
	if defaultNsResult.Err() == mongo.ErrNoDocuments {
		namespace = models.Namespace{
			Name:      "default",
			Workflows: make([]string, 0),
			Users:     []models.NamespaceUserRef{},
		}
		nsUser := models.NamespaceUserRef{
			Id:       admin.Id,
			Username: admin.Name,
			Accepted: true,
		}
		namespace.Users = append(namespace.Users, nsUser)
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
