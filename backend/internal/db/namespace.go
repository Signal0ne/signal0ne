package db

import (
	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNamespace(ctx *gin.Context, collection *mongo.Collection, filter bson.M) (models.Namespace, error) {
	var namespace models.Namespace
	err := collection.FindOne(ctx, filter).Decode(&namespace)
	if err != nil {
		return namespace, err
	}
	return namespace, nil
}

func AddUserToNamespace(ctx *gin.Context,
	collection *mongo.Collection,
	filter bson.M,
	user models.User) error {

	var userRef models.NamespaceUserRef

	userRef.Id = user.Id
	userRef.Username = user.Name
	userRef.Accepted = false

	_, err := collection.UpdateOne(ctx, filter, bson.M{"$push": bson.M{"users": userRef}})
	if err != nil {
		return err
	}
	return nil
}

func UpdateNamespace(ctx *gin.Context, collection *mongo.Collection, filter bson.M, update bson.M) error {
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
