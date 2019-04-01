package middle

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/krug-serv/api"
	"github.com/hugobrains/krug-serv/ir"
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

	/*
			pCache := bytes.NewBuffer(krugReq.Data)
		decCache := gob.NewDecoder(pCache)
		decCache.Decode(&payload)
	*/

	typedMod, errs := declType(payload.ScopeMap, payload.Module)

	jsonTypedMod, err := jsoniter.MarshalIndent(typedMod, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonTypedMod),
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
	/*
			pCache := bytes.NewBuffer(krugReq.Data)
		decCache := gob.NewDecoder(pCache)
		decCache.Decode(&irMod)
	*/

	scopeMap, errs := buildScope(irMod)

	jsonScopeMap, err := jsoniter.MarshalIndent(scopeMap, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonScopeMap),
		Errors: errs,
	}

	c.JSON(200, &resp)
}
