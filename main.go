package main

import (
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/krugc-api/front"
)

func main() {
	router := gin.Default()

	f := router.Group("/front")
	{
		f.POST("/lex", front.Tokenize)
	}

	router.Run()
}
