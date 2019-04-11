package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
	jsoniter "github.com/json-iterator/go"
)

func TypeResolve(c *gin.Context) {
	var krugReq api.TypeResolveRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	typeMap, errors := typeResolve(irMod)

	jsonIrModule, err := jsoniter.MarshalIndent(typeMap, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonIrModule),
		Errors: errors,
	}
	c.JSON(200, &resp)
}

func SymbolResolve(c *gin.Context) {
	var krugReq api.SymbolResolveRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	mod, errors := symResolve(irMod)

	jsonMod, err := jsoniter.MarshalIndent(mod, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonMod),
		Errors: errors,
	}
	c.JSON(200, &resp)
}
