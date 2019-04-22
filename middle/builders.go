package middle

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
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
	var buildTypeMapReq api.BuildTypeMapRequest
	if err := c.BindJSON(&buildTypeMapReq); err != nil {
		panic(err)
	}

	var payload struct {
		ScopeMap *ir.ScopeMap
		Module   *ir.Module
	}

	typedMod, errs := declType(payload.ScopeMap, payload.Module)

	jsonTypedMod, err := jsoniter.MarshalIndent(typedMod, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonTypedMod),
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}

// this returns the ir module, modified
// with the symbol tables. I feel like
// this should, however, just return the
// stab tree structure?
func BuildScope(c *gin.Context) {
	var scopeMapReq api.BuildScopeMapRequest
	if err := c.BindJSON(&scopeMapReq); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(scopeMapReq.IRModule), &irMod); err != nil {
		panic(err)
	}

	scopeMap, errs := buildScope(&irMod)

	jsonScopeMap, err := jsoniter.MarshalIndent(scopeMap, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonScopeMap),
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}

func BuildScopeDict(c *gin.Context) {
	var scopeDictReq api.BuildScopeDictRequest
	if err := c.BindJSON(&scopeDictReq); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(scopeDictReq.IRModule), &irMod); err != nil {
		panic(err)
	}

	scopeDict, errs := buildScopeDict(&irMod)

	jsonScopeDict, err := jsoniter.MarshalIndent(scopeDict, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonScopeDict),
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}
