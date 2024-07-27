package controllers

import "github.com/gin-gonic/gin"

type MainController struct {
}

func NewMainController() *MainController {
	return &MainController{}
}

func (c *MainController) ReceiveAlert(ctx *gin.Context) {

}

func (c *MainController) ApplyWorkflow(ctx *gin.Context) {

}
