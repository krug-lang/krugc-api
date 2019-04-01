package middle

import (
	"bytes"
	"encoding/gob"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/server/api"
	"github.com/krug-lang/server/ir"
)

func TypeResolve(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&irMod)

	typeMap, errors := typeResolve(irMod)
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	if err := encoder.Encode(&typeMap); err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
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
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&irMod)

	mod, errors := symResolve(irMod)

	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	if err := encoder.Encode(&mod); err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
		Errors: errors,
	}
	c.JSON(200, &resp)
}
