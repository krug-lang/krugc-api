package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/hugobrains/krug-serv/api"
	"github.com/hugobrains/krug-serv/ir"
	jsoniter "github.com/json-iterator/go"
)

func TypeResolve(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	/*
			pCache := bytes.NewBuffer(krugReq.Data)
		decCache := gob.NewDecoder(pCache)
		decCache.Decode(&irMod)
	*/

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
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	/*
			pCache := bytes.NewBuffer(krugReq.Data)
		decCache := gob.NewDecoder(pCache)
		decCache.Decode(&irMod)
	*/

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
