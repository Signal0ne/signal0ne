package middlewares

import (
	"net/http"
	"signal0ne/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	var jwtToken = strings.TrimPrefix(authHeader, "Bearer ")

	sourceHeader := ctx.GetHeader("X-Source")

	if sourceHeader == "" || sourceHeader == "frontend" {
		_, err := utils.VerifyToken(jwtToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
	} else if sourceHeader == "integration" {
		err := utils.VerifyIntegrationToken(jwtToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid source"})
		ctx.Abort()
		return
	}

	ctx.Next()
}
