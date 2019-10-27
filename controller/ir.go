package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/caasper/service"
)

type irService struct {}

func (is *irService) register(router *gin.Engine) {
	// intermediate representation, handles
	// conversion of AST into an IR
	i := router.Group("/ir")
	{
		i.POST("/build", service.Build)
	}
}
