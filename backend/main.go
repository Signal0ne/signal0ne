package main

import (
	"fmt"
	config "signal0ne/cmd/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var server = gin.Default()

	cfg := config.GetInstance()
	if cfg == nil {
		panic("CRITICAL: unable to load config")
	}

	fmt.Println("Hello, Signal0ne!")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.Run(":" + cfg.ServerPort)
}
