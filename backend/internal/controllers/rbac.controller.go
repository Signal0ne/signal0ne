package controllers

import (
	"net/http"
	"signal0ne/internal/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JoinRequest struct {
	NamespaceId string `json:"namespaceId"`
	UserId      string `json:"userId"`
}

type JoinTriageRequest struct {
	NamespaceId string `json:"namespaceId"`
	UserId      string `json:"userId"`
	Accepted    bool   `json:"accepted"`
}

type RBACController struct {
	UsersCollection      *mongo.Collection
	NamespacesCollection *mongo.Collection
}

func NewRBACController(
	usersCollection *mongo.Collection,
	namespacesCollection *mongo.Collection,
) *RBACController {
	return &RBACController{
		UsersCollection:      usersCollection,
		NamespacesCollection: namespacesCollection,
	}
}

func (c *RBACController) RequestToJoinNamespace(ctx *gin.Context) {
	var joinRequest JoinRequest
	err := ctx.ShouldBindJSON(&joinRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := primitive.ObjectIDFromHex(joinRequest.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUser(ctx, c.UsersCollection, bson.M{"_id": userId})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespaceId, err := primitive.ObjectIDFromHex(joinRequest.NamespaceId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace, err := db.GetNamespace(ctx, c.NamespacesCollection, bson.M{"_id": namespaceId})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, user := range namespace.Users {
		if user.Id == userId {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already in namespace"})
			return
		}
	}

	err = db.AddUserToNamespace(ctx, c.NamespacesCollection, bson.M{"_id": namespaceId}, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User added to namespace"})
}

func (c *RBACController) TriageNamespaceJoinRequest(ctx *gin.Context) {
	var joinTriageRequest JoinTriageRequest
	err := ctx.ShouldBindJSON(&joinTriageRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := primitive.ObjectIDFromHex(joinTriageRequest.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.GetUser(ctx, c.UsersCollection, bson.M{"_id": userId})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespaceId, err := primitive.ObjectIDFromHex(joinTriageRequest.NamespaceId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace, err := db.GetNamespace(ctx, c.NamespacesCollection, bson.M{"_id": namespaceId})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, user := range namespace.Users {
		if user.Id == userId {
			namespace.Users[i].Accepted = joinTriageRequest.Accepted
		}
	}

	err = db.UpdateNamespace(ctx, c.NamespacesCollection, bson.M{"_id": namespaceId}, bson.M{"$set": bson.M{"users": namespace.Users}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User triaged"})
}
