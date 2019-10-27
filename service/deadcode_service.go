package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/entity"
	"github.com/krug-lang/caasper/ir"
	"github.com/krug-lang/caasper/middle"
	"net/http"
)

// UnusedFunctions is a request that will check for
// any unused functions in the given module.
func UnusedFunctions(c *gin.Context) {
	var req entity.UnusedFunctionRequest
	if err := c.BindJSON(&req); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(req.IRModule), &irMod); err != nil {
		panic(err)
	}

	var scopeDict ir.ScopeDict
	if err := jsoniter.Unmarshal([]byte(req.ScopeMap), &scopeDict); err != nil {
		panic(err)
	}

	errs := middle.UnusedFunc(&irMod, &scopeDict)

	resp := entity.KrugResponse{
		Data:   "",
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}

