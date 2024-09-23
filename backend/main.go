package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"signal0ne/api/routers"
	"signal0ne/cmd/config"
	"signal0ne/internal/controllers"
	"signal0ne/internal/tools"
	"signal0ne/internal/utils"
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	alertsCollection := mongoConn.Database("signalone").Collection("alerts")
	namespacesCollection := mongoConn.Database("signalone").Collection("namespaces")
	workflowsCollection := mongoConn.Database("signalone").Collection("workflows")
	integrationsCollection := mongoConn.Database("signalone").Collection("integrations")
	incidentsCollection := mongoConn.Database("signalone").Collection("incidents")

	var conn net.Conn = nil
	conn = utils.ConnectToSocket()

	var logger *log.Logger = nil
	if cfg.Debug {
		logFile, err := os.OpenFile("/logs/workflow.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		logger = log.New(logFile, "", 0)
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
	namespaceController := controllers.NewNamespaceController(namespacesCollection)
	workflowController := controllers.NewWorkflowController(
		workflowsCollection,
		namespacesCollection,
		integrationsCollection,
		incidentsCollection,
		alertsCollection,
		cfg.Server,
		conn,
		logger,
	)
	integrationsController := controllers.NewIntegrationController(
		integrationsCollection,
		namespacesCollection,
	)
	incidentController := controllers.NewIncidentController(
		incidentsCollection,
		integrationsCollection,
		alertsCollection,
		workflowsCollection,
		namespacesCollection,
		conn,
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

	pyInterfacePayload := map[string]any{
		"command": "ping",
		"params":  map[string]any{},
	}

	payloadBytes, err := json.Marshal(pyInterfacePayload)
	if err != nil {
		panic(err)
	}

	if conn != nil {
		defer conn.Close()

		batchSizeHeader := make([]byte, 4)
		binary.BigEndian.PutUint32(batchSizeHeader, uint32(len(payloadBytes)))
		payloadBytesWithHeaders := append(batchSizeHeader, payloadBytes...)

		_, err = conn.Write(payloadBytesWithHeaders)
		if err != nil {
			panic(err)
		}

		headerBuffer := make([]byte, 4)

		_, err = conn.Read(headerBuffer)
		if err != nil {
			panic(err)
		}

		size := binary.BigEndian.Uint32(headerBuffer)

		payloadBuffer := make([]byte, size)
		_, err = conn.Read(payloadBuffer)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("---------------------- WARNING! ----------------------------")
		fmt.Println("Failed to establish connection to python interface.\n It is possible that the interface is not running. \n Enable it to use the full functionality of the application.")
		fmt.Println("---------------------- WARNING! ----------------------------")
	}

	server.Run(":" + cfg.Server.ServerPort)
}
