package middle

import (
	"bytes"
	"encoding/gob"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/krugc-api/api"
	"github.com/krug-lang/krugc-api/ir"
)

/*
	middle end:

	declare types exist (ir builder does this)

	check all references to types are OK
		check structure bodies.
		check function bodies.
		check function params
*/

func TypeResolve(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&irMod)

	errors := typeResolve(irMod)

	resp := api.KrugResponse{
		Data:   []byte{},
		Errors: errors,
	}
	c.JSON(200, &resp)
}

// this returns the ir module, modified
// with the symbol tables. I feel like
// this should, however, just return the
// stab tree structure?
func BuildScope(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&irMod)

	scopedMod, errs := buildScope(irMod)

	// write new module with built scopes.
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	if err := encoder.Encode(&scopedMod); err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
		Errors: errs,
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

	//TODO

	resp := api.KrugResponse{
		Data:   []byte{},
		Errors: nil,
	}
	c.JSON(200, &resp)
}
