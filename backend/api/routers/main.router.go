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
}

func NewMainRouter(
	MainController *controllers.MainController,
	NamespaceController *controllers.NamespaceController,
	WorkflowController *controllers.WorkflowController,
	IntegrationController *controllers.IntegrationController,
	UserAuthController *controllers.UserAuthController,
) *MainRouter {
	return &MainRouter{
		IntegrationController: IntegrationController,
		MainController:        MainController,
		NamespaceController:   NamespaceController,
		UserAuthController:    UserAuthController,
		WorkflowController:    WorkflowController,
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
		workflowGroup.GET("/:workflowid/get")
		workflowGroup.DELETE("/:workflowid/delete")
		workflowGroup.PATCH("/:workflowid/update")
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
