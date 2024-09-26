package controllers

import "github.com/gin-gonic/gin"

type MainController struct {
}

func NewMainController() *MainController {
	return &MainController{}
}

func (c *MainController) HealthCheck(ctx *gin.Context) {

}
