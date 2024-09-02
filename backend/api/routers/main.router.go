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
	IncidentController    *controllers.IncidentController
}

func NewMainRouter(
	MainController *controllers.MainController,
	NamespaceController *controllers.NamespaceController,
	WorkflowController *controllers.WorkflowController,
	IntegrationController *controllers.IntegrationController,
	IncidentController *controllers.IncidentController,
	UserAuthController *controllers.UserAuthController,
) *MainRouter {
	return &MainRouter{
		IntegrationController: IntegrationController,
		MainController:        MainController,
		NamespaceController:   NamespaceController,
		UserAuthController:    UserAuthController,
		WorkflowController:    WorkflowController,
		IncidentController:    IncidentController,
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

	incidentGroup := rg.Group("/:namespaceid/incident")
	{
		incidentGroup.GET("/:incidentid")
	}

	integrationGroup := rg.Group("/:namespaceid/integration", middlewares.CheckAuthorization)
	{
		integrationGroup.POST("/create", r.IntegrationController.Install)
		integrationGroup.GET("/installable", r.IntegrationController.GetInstallableIntegrations)
		integrationGroup.GET("/installed", r.IntegrationController.GetInstalledIntegrations)
		integrationGroup.DELETE("/:integrationid/delete")
		integrationGroup.GET("/:integrationid/get")
		integrationGroup.PATCH("/:integrationid/update")
	}

	namespaceGroup := rg.Group("/namespace")
	{
		namespaceGroup.GET("/search-by-name", r.NamespaceController.GetNamespaceByName)
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
		workflowGroup.GET("/workflows", r.WorkflowController.GetWorkflows)
		workflowGroup.PATCH("/:workflowid")
	}
}
