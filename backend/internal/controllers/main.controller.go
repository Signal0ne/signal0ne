package controllers

import "github.com/gin-gonic/gin"

type MainController struct {
}

func NewMainController() *MainController {
	return &MainController{}
}

func (c *NamespaceController) HealthCheck(ctx *gin.Context) {

}
