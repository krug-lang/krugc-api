package middle

import (
	"bytes"
	"encoding/gob"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/server/api"
	"github.com/krug-lang/server/ir"
)

/*
	middle end:

	declare types exist (ir builder does this)

	check all references to types are OK
		check structure bodies.
		check function bodies.
		check function params
*/

func BuildType(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var payload struct {
		ScopeMap *ir.ScopeMap
		Module   *ir.Module
	}

	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&payload)

	typedMod, errs := declType(payload.ScopeMap, payload.Module)

	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	if err := encoder.Encode(&typedMod); err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
		Errors: errs,
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

	scopeMap, errs := buildScope(irMod)

	// write new module with built scopes.
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	if err := encoder.Encode(&scopeMap); err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
		Errors: errs,
	}

	c.JSON(200, &resp)
}
