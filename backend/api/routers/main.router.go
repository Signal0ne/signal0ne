package routers

import (
	"signal0ne/internal/controllers"
	"signal0ne/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type MainRouter struct {
	AlertController       *controllers.AlertController
	IncidentController    *controllers.IncidentController
	IntegrationController *controllers.IntegrationController
	MainController        *controllers.MainController
	NamespaceController   *controllers.NamespaceController
	RBACController        *controllers.RBACController
	UserAuthController    *controllers.UserAuthController
	WorkflowController    *controllers.WorkflowController
}

func NewMainRouter(
	AlertController *controllers.AlertController,
	IncidentController *controllers.IncidentController,
	IntegrationController *controllers.IntegrationController,
	MainController *controllers.MainController,
	NamespaceController *controllers.NamespaceController,
	RBACController *controllers.RBACController,
	UserAuthController *controllers.UserAuthController,
	WorkflowController *controllers.WorkflowController,
) *MainRouter {
	return &MainRouter{
		AlertController:       AlertController,
		IncidentController:    IncidentController,
		IntegrationController: IntegrationController,
		MainController:        MainController,
		NamespaceController:   NamespaceController,
		RBACController:        RBACController,
		UserAuthController:    UserAuthController,
		WorkflowController:    WorkflowController,
	}
}

func (r *MainRouter) RegisterRoutes(rg *gin.RouterGroup) {

	alertGroup := rg.Group("/:namespaceid/alert", middlewares.CheckAuthorization)
	{
		alertGroup.GET("/:alertid", r.AlertController.GetAlert)
	}

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", r.UserAuthController.Login)
		authGroup.GET("/logout", r.UserAuthController.Logout)
		authGroup.POST("/register", r.UserAuthController.Register)
		authGroup.GET("/token/refresh", r.UserAuthController.RefreshToken)
	}

	rbacGroup := rg.Group("/rbac", middlewares.CheckAuthorization)
	{
		rbacGroup.POST("/request-namespace-access", r.RBACController.RequestToJoinNamespace)
		rbacGroup.POST("/namespace-access-response", r.RBACController.TriageNamespaceJoinRequest)
	}

	incidentGroup := rg.Group("/:namespaceid/incident", middlewares.CheckAuthorization)
	{
		// All incident types
		incidentGroup.POST("", r.IncidentController.CreateIncident)
		incidentGroup.GET("incidents", r.IncidentController.GetIncidents)
		incidentGroup.POST("/:incidentid/tasks", r.IncidentController.AddNewTask)
		incidentGroup.POST("/:incidentid/:taskid/add-task-comment", r.IncidentController.AddTaskComment)
		incidentGroup.PATCH("/:incidentid/:taskid/assignee", r.IncidentController.UpdateTaskAssignee)
		incidentGroup.PATCH("/:incidentid/:taskid/status", r.IncidentController.UpdateTaskStatus)
		incidentGroup.PATCH("/:incidentid/update-tasks-priority", r.IncidentController.UpdateTasksPriority)

		// Signal0ne incident only
		incidentGroup.PATCH("/:incidentid", r.IncidentController.UpdateIncident)
		incidentGroup.GET("/:incidentid", r.IncidentController.GetIncident)
		incidentGroup.POST("/:incidentid/register-history-event/:updatetype", r.IncidentController.RegisterHistoryEvent)
	}

	integrationGroup := rg.Group("/:namespaceid/integration", middlewares.CheckAuthorization)
	{
		integrationGroup.POST("", r.IntegrationController.Install)
		integrationGroup.GET("/installable", r.IntegrationController.GetInstallableIntegrations)
		integrationGroup.GET("/installed", r.IntegrationController.GetInstalledIntegrations)
		integrationGroup.DELETE("/:integrationid")
		integrationGroup.GET("/:integrationid", r.IntegrationController.GetIntegration)
		integrationGroup.PATCH("/:integrationid", r.IntegrationController.UpdateIntegration)
	}

	namespaceGroup := rg.Group("/namespace", middlewares.CheckAuthorization)
	{
		namespaceGroup.POST("/create")
		namespaceGroup.DELETE("/:namespaceid/delete")
		namespaceGroup.GET("/:namespaceid/get")
		namespaceGroup.PATCH("/:namespaceid/update")
		namespaceGroup.GET("/:namespaceid/users", r.NamespaceController.GetUsersFromNamespace)
		namespaceGroup.GET("/search-by-name", r.NamespaceController.GetNamespaceByName)
	}

	webhookGroup := rg.Group("/webhook")
	{
		webhookGroup.POST("/:namespaceid/:workflowid/:salt", r.WorkflowController.WebhookTriggerHandler)
	}

	workflowGroup := rg.Group("/:namespaceid/workflow", middlewares.CheckAuthorization)
	{
		workflowGroup.POST("/create", r.WorkflowController.ApplyWorkflow)
		workflowGroup.DELETE("/:workflowid")
		workflowGroup.GET("/:workflowid", r.WorkflowController.GetWorkflow)
		workflowGroup.GET("/workflows", r.WorkflowController.GetWorkflows)
		workflowGroup.PATCH("/:workflowid")
	}
}
