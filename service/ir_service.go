package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/entity"
	"github.com/krug-lang/caasper/front"
	"github.com/krug-lang/caasper/ir"
	"net/http"
)

// Build takes the given []ParseTree's
// and builds a SINGLE ir module from them.
func Build(c *gin.Context) {
	var irBuildReq entity.IRBuildRequest
	if err := c.BindJSON(&irBuildReq); err != nil {
		panic(err)
	}

	var trees [][]*front.ParseTreeNode

	if err := jsoniter.Unmarshal([]byte(irBuildReq.TreeNodes), &trees); err != nil {
		panic(err)
	}

	irModule, errors := ir.Build(trees)

	jsonIrModule, err := jsoniter.MarshalIndent(irModule, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonIrModule),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}

