package controllers

import (
	"fmt"
	"net/http"
	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NamespaceController struct {
	NamespaceCollection *mongo.Collection

	// ==== Use as readonly ====
	UsersCollection *mongo.Collection
	// =========================
}

func NewNamespaceController(namespaceCollection *mongo.Collection, usersCollection *mongo.Collection) *NamespaceController {
	return &NamespaceController{
		NamespaceCollection: namespaceCollection,
		UsersCollection:     usersCollection,
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

func (c *NamespaceController) GetUsersFromNamespace(ctx *gin.Context) {
	var finalUsers []models.User
	var namespace *models.Namespace
	var userIds []primitive.ObjectID

	namespaceId := ctx.Param("namespaceid")

	nsID, err := primitive.ObjectIDFromHex(namespaceId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid namespace ID: %v", err),
		})

		return
	}

	res := c.NamespaceCollection.FindOne(ctx, primitive.M{"_id": nsID})
	err = res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot find namespace: %v", err),
		})
		return
	}

	for _, user := range namespace.Users {
		if user.Accepted {
			userId, err := primitive.ObjectIDFromHex(user.Id.Hex())
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Invalid user ID: %v", err),
				})

				return
			}

			userIds = append(userIds, userId)
		}
	}

	filter := bson.M{"_id": bson.M{"$in": userIds}}

	projection := bson.M{"password": 0}

	cursor, err := c.UsersCollection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error fetching users: %v", err),
		})

		return
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &finalUsers); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error decoding users: %v", err),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": finalUsers,
	})
}
