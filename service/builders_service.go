package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/entity"
	"github.com/krug-lang/caasper/ir"
	"github.com/krug-lang/caasper/middle"
	"net/http"
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
	var buildTypeMapReq entity.BuildTypeMapRequest
	if err := c.BindJSON(&buildTypeMapReq); err != nil {
		panic(err)
	}

	var payload struct {
		ScopeMap *ir.ScopeMap
		Module   *ir.Module
	}

	typedMod, errs := middle.DeclType(payload.ScopeMap, payload.Module)

	jsonTypedMod, err := jsoniter.MarshalIndent(typedMod, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
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
	var scopeMapReq entity.BuildScopeMapRequest
	if err := c.BindJSON(&scopeMapReq); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(scopeMapReq.IRModule), &irMod); err != nil {
		panic(err)
	}

	scopeMap, errs := middle.BuildScope(&irMod)

	jsonScopeMap, err := jsoniter.MarshalIndent(scopeMap, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonScopeMap),
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}

func BuildScopeDict(c *gin.Context) {
	var scopeDictReq entity.BuildScopeDictRequest
	if err := c.BindJSON(&scopeDictReq); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(scopeDictReq.IRModule), &irMod); err != nil {
		panic(err)
	}

	scopeDict, errs := middle.BuildScopeDict(&irMod)

	jsonScopeDict, err := jsoniter.MarshalIndent(scopeDict, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonScopeDict),
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}

