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

	namespaceGroup := rg.Group("/namespace")
	{
		namespaceGroup.POST("/create")
		namespaceGroup.GET("/:namespaceid/get")
		namespaceGroup.DELETE("/:namespaceid/delete")
		namespaceGroup.PATCH("/:namespaceid/update")
	}

	workflowGroup := rg.Group("/:namespaceid/workflow")
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

	integrationGroup := rg.Group("/:namespaceid/integration")
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
