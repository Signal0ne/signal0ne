package routers

import (
	"signal0ne/internal/controllers"

	"github.com/gin-gonic/gin"
)

type MainRouter struct {
	Controller *controllers.MainController
}

func NewMainRouter(MainController *controllers.MainController) *MainRouter {
	return &MainRouter{
		Controller: MainController,
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
		workflowGroup.POST("/create", r.Controller.ApplyWorkflow)
		workflowGroup.GET("/:workflowid/get")
		workflowGroup.DELETE("/:workflowid/delete")
		workflowGroup.PATCH("/:workflowid/update")
	}

	webhookGroup := rg.Group("/webhook")
	{
		webhookGroup.POST("/:namespaceid/:workflowid/:salt", r.Controller.ReceiveAlert)
	}
}
