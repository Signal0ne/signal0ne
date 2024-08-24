package middlewares

import (
	"fmt"
	"net/http"
	"signal0ne/cmd/config"
	"signal0ne/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(ctx *gin.Context) {
	var skipAuth = config.GetInstance().SkipAuth

	//TODO: Remove before release
	fmt.Printf("SKIP AUTH: %v\n", skipAuth)

	if skipAuth {
		ctx.Next()
		return
	}

	authHeader := ctx.GetHeader("Authorization")

	var jwtToken = strings.TrimPrefix(authHeader, "Bearer ")

	_, err := utils.VerifyToken(jwtToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Next()
}
