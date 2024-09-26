package controllers

import (
	"net/http"
	"signal0ne/internal/db"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UserAuthController struct {
	UsersCollection      *mongo.Collection
	NamespacesCollection *mongo.Collection
}

func NewUserAuthController(
	usersCollection *mongo.Collection,
	namespacesCollection *mongo.Collection,
) *UserAuthController {
	return &UserAuthController{
		UsersCollection:      usersCollection,
		NamespacesCollection: namespacesCollection,
	}
}

func (c *UserAuthController) Register(ctx *gin.Context) {
	var registerUserRequest UserAuthRequest
	var user models.User

	err := ctx.ShouldBindJSON(&registerUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(registerUserRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify if admin user exists
	filter := bson.M{
		"role": "admin",
	}
	_, err = db.GetUser(ctx, c.UsersCollection, filter)
	if err == mongo.ErrNoDocuments {

		user.Id = primitive.NewObjectID()
		user.Name = registerUserRequest.Username
		user.Password = hashedPassword
		user.Role = models.AdminRole

		err = tools.OnboardAdmin(ctx, c.NamespacesCollection, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if err != nil && err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {

		user.Id = primitive.NewObjectID()
		user.Name = registerUserRequest.Username
		user.Password = hashedPassword
		user.Role = models.UserRole
	}

	err = db.CreateUser(ctx, c.UsersCollection, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.CreateToken(user, "access")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(user, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (c *UserAuthController) Login(ctx *gin.Context) {
	var registerUserRequest UserAuthRequest
	var user models.User

	err := ctx.ShouldBindJSON(&registerUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err = db.GetUser(ctx, c.UsersCollection, bson.M{"name": registerUserRequest.Username})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	match := utils.ComparePasswordHashes(user.Password, registerUserRequest.Password)
	if !match {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := utils.CreateToken(user, "access")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(user, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (c *UserAuthController) RefreshToken(ctx *gin.Context) {
	var data RefreshTokenRequest
	var user models.User

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := utils.VerifyToken(data.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err = db.GetUser(ctx, c.UsersCollection, bson.M{"_id": id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.CreateToken(user, "access")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(user, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
