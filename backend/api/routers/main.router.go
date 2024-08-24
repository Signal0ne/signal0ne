package routers

import (
	"signal0ne/internal/controllers"
	"signal0ne/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type MainRouter struct {
	IntegrationController *controllers.IntegrationController
	MainController        *controllers.MainController
	NamespaceController   *controllers.NamespaceController
	UserAuthController    *controllers.UserAuthController
	WorkflowController    *controllers.WorkflowController
	AlertController       *controllers.AlertsController
}

func NewMainRouter(
	MainController *controllers.MainController,
	NamespaceController *controllers.NamespaceController,
	WorkflowController *controllers.WorkflowController,
	IntegrationController *controllers.IntegrationController,
	AlertController *controllers.AlertsController,
	UserAuthController *controllers.UserAuthController,
) *MainRouter {
	return &MainRouter{
		IntegrationController: IntegrationController,
		MainController:        MainController,
		NamespaceController:   NamespaceController,
		UserAuthController:    UserAuthController,
		WorkflowController:    WorkflowController,
		IntegrationController: IntegrationController,
		AlertController:       AlertController,
	}
}

func (r *MainRouter) RegisterRoutes(rg *gin.RouterGroup) {

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/email-confirmation")
		authGroup.POST("/email-confirmation-link-resend")
		authGroup.POST("/login")
		authGroup.POST("/register")
		authGroup.POST("/token/refresh")
	}

	namespaceGroup := rg.Group("/namespace", middlewares.CheckAuthorization)
	{
		namespaceGroup.POST("/create")
		namespaceGroup.GET("/:namespaceid/get")
		namespaceGroup.DELETE("/:namespaceid/delete")
		namespaceGroup.PATCH("/:namespaceid/update")
	}

	workflowGroup := rg.Group("/:namespaceid/workflow", middlewares.CheckAuthorization)
	{
		workflowGroup.POST("/create", r.WorkflowController.ApplyWorkflow)
		workflowGroup.GET("/:workflowid")
		workflowGroup.DELETE("/:workflowid")
		workflowGroup.PATCH("/:workflowid")
	}

	alertGroup := rg.Group("/alert")
	{
		alertGroup.GET("/:alertid/details", r.AlertController.Details)
		alertGroup.GET("/:alertid/summary", r.AlertController.Summary)
		alertGroup.GET("/:alertid/correlations", r.AlertController.Correlations)
	}

	integrationGroup := rg.Group("/:namespaceid/integration", middlewares.CheckAuthorization)
	{
		integrationGroup.POST("/create", r.IntegrationController.Install)
		integrationGroup.GET("/:integrationid/get")
		integrationGroup.DELETE("/:integrationid/delete")
		integrationGroup.PATCH("/:integrationid/update")
	}

	webhookGroup := rg.Group("/webhook")
	{
		webhookGroup.POST("/:namespaceid/:workflowid/:salt", r.WorkflowController.WebhookTriggerHandler)
	}
}
