package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/caasper/service"
)

type frontendService struct {}

func (fs *frontendService) register(router *gin.Engine) {
	// compiler frontend, handles lexing/parsing
	f := router.Group("/front")
	{
		// lexical analysis
		f.POST("/lex", service.Tokenize)

		// parsing.
		parse := f.Group("/parse")
		{
			parse.POST("/ast", service.Parse)
			parse.POST("/directive", service.DirectiveParser)
		}

		f.POST("/comments", service.Comments)
	}
}

