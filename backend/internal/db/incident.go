package db

import (
	"context"
	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetIncidentById(id string,
	ctx context.Context,
	incidentsCollection *mongo.Collection,
) (models.Incident, error) {
	var incident models.Incident

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return incident, err
	}

	err = incidentsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&incident)
	if err != nil {
		return incident, err
	}

	return incident, nil
}

func SetRegisterHistoryEvent(incidentId string,
	incidentUpdate any,
	ctx *gin.Context,
	incidentsCollection *mongo.Collection,
) error {
	objectID, err := primitive.ObjectIDFromHex(incidentId)
	if err != nil {
		return err
	}
	_, err = incidentsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$push": bson.M{"history": incidentUpdate}},
	)
	if err != nil {
		return err
	}

	return nil
}
