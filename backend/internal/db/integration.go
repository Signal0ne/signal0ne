package db

import (
	"context"
	"signal0ne/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetIntegrationByName(
	integrationName string,
	ctx context.Context,
	integrationsCollection *mongo.Collection,
) (models.Integration, error) {
	res := integrationsCollection.FindOne(ctx, bson.M{"name": integrationName})

	var integration models.Integration
	err := res.Decode(&integration)
	if err != nil {
		return models.Integration{}, err
	}

	return integration, nil
}
