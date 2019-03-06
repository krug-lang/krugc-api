package front

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Tokenize(c *gin.Context) {
	fmt.Println("tokenize time!")

	var sourceFile KrugCompilationUnit
	if err := c.BindJSON(&sourceFile); err != nil {
		panic(err)
	}

	tokens := []*Token{}

	fmt.Println("Tokenizing ", sourceFile)
	c.JSON(200, TokenStream{tokens})
}
