package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/caasper/back"
)

type backendService struct {}

func (bs *backendService) register(router *gin.Engine) {
	// backend of the compiler handles taking the
	// IR and generating code from it.
	b := router.Group("/back")
	{
		b.POST("/gen", back.Gen)
	}
}
