package routers

import (
	"signal0ne/internal/controllers"

	"github.com/gin-gonic/gin"
)

type MainRouter struct {
	MainController        *controllers.MainController
	NamespaceController   *controllers.NamespaceController
	WorkflowController    *controllers.WorkflowController
	IntegrationController *controllers.IntegrationController
	AlertController       *controllers.AlertsController
}

func NewMainRouter(
	MainController *controllers.MainController,
	NamespaceController *controllers.NamespaceController,
	WorkflowController *controllers.WorkflowController,
	IntegrationController *controllers.IntegrationController,
	AlertController *controllers.AlertsController,
) *MainRouter {
	return &MainRouter{
		MainController:        MainController,
		NamespaceController:   NamespaceController,
		WorkflowController:    WorkflowController,
		IntegrationController: IntegrationController,
		AlertController:       AlertController,
	}
}

func (r *MainRouter) RegisterRoutes(rg *gin.RouterGroup) {

	alertGroup := rg.Group("/alert")
	{
		alertGroup.GET("/:alertid/correlations", r.AlertController.Correlations)
		alertGroup.GET("/:alertid/details", r.AlertController.Details)
		alertGroup.GET("/:alertid/summary", r.AlertController.Summary)
	}

	integrationGroup := rg.Group("/:namespaceid/integration")
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
