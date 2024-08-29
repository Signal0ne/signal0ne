package main

import (
	"context"
	"fmt"
	"net"
	"signal0ne/api/routers"
	"signal0ne/cmd/config"
	"signal0ne/internal/controllers"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations"
	"strings"
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoConn, err := tools.InitMongoClient(ctx, cfg.MongoUri)
	if err != nil {
		panic(
			fmt.Sprintf("Failed to establish connection to %s, error: %s",
				strings.Split(cfg.MongoUri, "/")[2],
				err),
		)
	}
	defer mongoConn.Disconnect(ctx)

	namespacesCollection := mongoConn.Database("signalone").Collection("namespaces")
	workflowsCollection := mongoConn.Database("signalone").Collection("workflows")
	integrationsCollection := mongoConn.Database("signalone").Collection("integrations")
	alertsCollection := mongoConn.Database("signalone").Collection("alerts")

	conn, err := net.DialTimeout("unix", cfg.IPCSocket, (15 * time.Second))
	if err != nil {
		panic(
			fmt.Sprintf("Failed to establish connection to %s, error: %s",
				cfg.IPCSocket,
				err),
		)
	} else {
		defer conn.Close()
	}

	err = tools.Initialize(ctx, namespacesCollection)
	if err != nil {
		panic(err)
	}

	// Loading installable integrations
	_, err = integrations.GetInstallableIntegrationsLib()
	if err != nil {
		panic(err)
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
	namespaceController := controllers.NewNamespaceController()
	workflowController := controllers.NewWorkflowController(
		workflowsCollection,
		namespacesCollection,
		integrationsCollection,
		alertsCollection,
		cfg.Server,
		conn)
	integrationsController := controllers.NewIntegrationController(
		integrationsCollection,
		namespacesCollection,
	)
	incidentController := controllers.NewIncidentController(
		alertsCollection,
	)
	userAuthController := controllers.NewUserAuthController()

	mainRouter := routers.NewMainRouter(
		mainController,
		namespaceController,
		workflowController,
		integrationsController,
		incidentController,
		userAuthController)
	mainRouter.RegisterRoutes(routerApiGroup)

	server.Run(":" + cfg.Server.ServerPort)
}
