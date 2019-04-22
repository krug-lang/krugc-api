package middle

/*
	not sure what to call this pass, but it will
	basically go through each function in a module
	and will check the move semantics are correct
	for owned memory/value bindings.
*/

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
	jsoniter "github.com/json-iterator/go"
)

type borrowChecker struct {
}

func (b *borrowChecker) visit(fn *ir.SymbolTable) {
	fmt.Println(fn.SymbolSet)
}

func (b *borrowChecker) validate(rootStab *ir.SymbolTable) {
	toVisit := []*ir.SymbolTable{}
	toVisit = append(toVisit, rootStab)

	for len(toVisit) != 0 {
		top := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]

		b.visit(top)

		for i := len(top.Inner) - 1; i >= 0; i-- {
			block := top.Inner[i]
			toVisit = append(toVisit, block)
		}
	}
}

func borrowCheck(mod *ir.Module, scopeMap *ir.ScopeMap) []api.CompilerError {
	errs := []api.CompilerError{}

	for name, fn := range scopeMap.Functions {
		checker := &borrowChecker{}
		fmt.Println("validating ", name)
		checker.validate(fn)
	}

	return errs
}

func BorrowCheck(c *gin.Context) {
	var req api.BorrowCheckRequest
	if err := c.BindJSON(&req); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(req.IRModule), &irMod); err != nil {
		panic(err)
	}

	var scopeMap ir.ScopeMap
	if err := jsoniter.Unmarshal([]byte(req.ScopeMap), &scopeMap); err != nil {
		panic(err)
	}

	errs := borrowCheck(&irMod, &scopeMap)

	resp := api.KrugResponse{
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}
