package main

import (
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/krugc-api/back"
	"github.com/krug-lang/krugc-api/front"
	"github.com/krug-lang/krugc-api/ir"
)

func main() {
	router := gin.Default()

	f := router.Group("/front")
	{
		// lexical analysis
		f.POST("/lex", front.Tokenize)

		// parsing.
		f.POST("/parse", front.Parse)
	}

	i := router.Group("/ir")
	{
		i.POST("/build", ir.Build)
	}

	b := router.Group("/back")
	{
		b.POST("/gen", back.Gen)
	}

	router.Run("localhost:8080")
}

func Now() {
	main()
}
