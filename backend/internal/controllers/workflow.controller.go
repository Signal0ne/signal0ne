package controllers

import "github.com/gin-gonic/gin"

type WorkflowController struct {
}

func NewWorkflowController() *WorkflowController {
	return &WorkflowController{}
}

func (c *WorkflowController) ReceiveAlert(ctx *gin.Context) {

}

func (c *WorkflowController) ApplyWorkflow(ctx *gin.Context) {

}
