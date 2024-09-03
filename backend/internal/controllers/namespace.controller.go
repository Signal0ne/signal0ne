package controllers

import (
	"fmt"
	"net/http"
	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NamespaceController struct {
	NamespaceCollection *mongo.Collection
}

func NewNamespaceController(namespaceCollection *mongo.Collection) *NamespaceController {
	return &NamespaceController{
		NamespaceCollection: namespaceCollection,
	}
}

func (c *NamespaceController) CreateOrUpdateNamespace(ctx *gin.Context) {

}

func (c *NamespaceController) GetNamespace(ctx *gin.Context) {

}

func (c *NamespaceController) GetNamespaceByName(ctx *gin.Context) {
	var namespace *models.Namespace

	namespaceName := ctx.Query("name")

	if namespaceName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	results := c.NamespaceCollection.FindOne(ctx, bson.M{"name": namespaceName})
	err := results.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("cannot find namespace, %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"namespaceId": namespace.Id})
}
