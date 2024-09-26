package db

import (
	"context"
	"signal0ne/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(
	ctx context.Context,
	usersCollection *mongo.Collection,
	filter bson.M,
) (models.User, error) {

	var user models.User

	res := usersCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return user, res.Err()
	}

	err := res.Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func CreateUser(
	ctx context.Context,
	usersCollection *mongo.Collection,
	user models.User,
) error {

	_, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
