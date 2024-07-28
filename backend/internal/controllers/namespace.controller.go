package controllers

import "github.com/gin-gonic/gin"

type NamespaceController struct {
}

func NewNamespaceController() *NamespaceController {
	return &NamespaceController{}
}

func (c *NamespaceController) CreateOrUpdateNamespace(ctx *gin.Context) {

}

func (c *NamespaceController) GetNamespace(ctx *gin.Context) {

}
