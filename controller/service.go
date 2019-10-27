package controller

import "github.com/gin-gonic/gin"

type endpointService interface {
	register(router *gin.Engine)
}
