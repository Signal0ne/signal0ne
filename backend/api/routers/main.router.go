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

	alertGroup := rg.Group("/alert")
	{
		alertGroup.GET("/:alertid/correlations", r.AlertController.Correlations)
		alertGroup.GET("/:alertid/details", r.AlertController.Details)
		alertGroup.GET("/:alertid/summary", r.AlertController.Summary)
	}

	integrationGroup := rg.Group("/:namespaceid/integration", middlewares.CheckAuthorization)
	{
		integrationGroup.POST("/create", r.IntegrationController.Install)
		integrationGroup.DELETE("/:integrationid/delete")
		integrationGroup.GET("/:integrationid/get")
		integrationGroup.PATCH("/:integrationid/update")
	}

	namespaceGroup := rg.Group("/namespace")
	{
		namespaceGroup.POST("/create")
		namespaceGroup.DELETE("/:namespaceid/delete")
		namespaceGroup.GET("/:namespaceid/get")
		namespaceGroup.PATCH("/:namespaceid/update")
	}

	webhookGroup := rg.Group("/webhook")
	{
		webhookGroup.POST("/:namespaceid/:workflowid/:salt", r.WorkflowController.WebhookTriggerHandler)
	}

	workflowGroup := rg.Group("/:namespaceid/workflow")
	{
		workflowGroup.POST("/create", r.WorkflowController.ApplyWorkflow)
		workflowGroup.DELETE("/:workflowid")
		workflowGroup.GET("/:workflowid")
		workflowGroup.PATCH("/:workflowid")
	}
}
