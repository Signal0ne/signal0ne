package main

import (
	"fmt"
	"net"
	"signal0ne/api/routers"
	"signal0ne/cmd/config"
	"signal0ne/internal/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var server = gin.Default()

	cfg := config.GetInstance()
	if cfg == nil {
		panic("CRITICAL: unable to load config")
	}

	socketPath := "/net/socket"

	conn, err := net.DialTimeout("unix", socketPath, (15 * time.Second))
	if err != nil {
		panic(fmt.Sprintf("Failed to establish connectiom to %s, error: %s", socketPath, err))
	} else {
		defer conn.Close()
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	routerApiGroup := server.Group("/api")
	routerApiGroup.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	mainController := controllers.NewMainController()
	mainRouter := routers.NewMainRouter(mainController)
	mainRouter.RegisterRoutes(routerApiGroup)

	//==========REMOVE BEFORE RELEASE==========
	_, err = conn.Write([]byte("Hello I am Go!"))
	if err != nil {
		fmt.Printf("Failed to send data: %s", err)
	}

	// Receive response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to read response: %s", err)
	}

	fmt.Printf("%s\n", buffer[:n])
	//===================

	server.Run(":" + cfg.ServerPort)
}
