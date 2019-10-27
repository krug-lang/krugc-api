package middle

import (
	"github.com/krug-lang/caasper/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/ir"
)

func TypeResolve(c *gin.Context) {
	var krugReq entity.TypeResolveRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	typeMap, errors := typeResolve(irMod)

	jsonIrModule, err := jsoniter.MarshalIndent(typeMap, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonIrModule),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}

func SymbolResolve(c *gin.Context) {
	var krugReq entity.SymbolResolveRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	mod, errors := symResolve(irMod)

	jsonMod, err := jsoniter.MarshalIndent(mod, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonMod),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}
