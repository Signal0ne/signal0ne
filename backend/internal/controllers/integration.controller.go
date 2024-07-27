package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type IntegrationController struct {
	IntegrationCollection *mongo.Collection
	NamespaceCollection   *mongo.Collection //Must be used as Readonly
}

func NewIntegrationController(
	integrationCollection *mongo.Collection,
	namespaceCollection *mongo.Collection) *IntegrationController {
	return &IntegrationController{
		IntegrationCollection: integrationCollection,
		NamespaceCollection:   namespaceCollection,
	}
}

func (ic *IntegrationController) CreateOrUpdate(ctx *gin.Context) {
}
