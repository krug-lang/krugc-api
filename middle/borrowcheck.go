package middle

/*
	not sure what to call this pass, but it will
	basically go through each function in a module
	and will check the move semantics are correct
	for owned memory/value bindings.

	TODO:
	move this into some kind of visitor
	so that we can reuse this over the sema pass
*/

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
	jsoniter "github.com/json-iterator/go"
)

type borrowChecker struct {
	block     *ir.Block
	depth     int
	scopeDict *ir.ScopeDict
}

/*
	we need the scope map to somehow map to the ir builder blocks...
	how could this work?
*/

func printSymbol(val *ir.SymbolValue) {
	switch val.Kind {
	case ir.SymbolKind:
		fmt.Println(val.Symbol)
	case ir.SymbolTableKind:
		fmt.Println(val.SymbolTable)
	}
}

func (b *borrowChecker) printSymScope(fn *ir.Function, rootStab *ir.SymbolTable) {
	var validateTree func(head *ir.SymbolTable, depth int)
	validateTree = func(head *ir.SymbolTable, depth int) {
		if head == nil {
			return
		}

		tab := strings.Repeat(" ", depth*4)
		fmt.Print(tab)
		fmt.Println(head.SymbolSet)

		for _, block := range head.Inner {
			validateTree(block, depth+1)
		}
	}
	validateTree(rootStab, 0)
}

func (b *borrowChecker) visitBlock(block *ir.Block) {
	for _, instr := range block.Instr {
		switch instr.Kind {
		case ir.BlockInstr:
			b.pushBlock(instr.Block)
		default:
			b.visitInstr(block, instr)
		}
	}
}

func (b *borrowChecker) pushBlock(block *ir.Block) {
	b.depth++
	fmt.Println("\npushed block", b.depth)
	b.visitBlock(block)
}

func (b *borrowChecker) visitInstr(parent *ir.Block, instr *ir.Instruction) {
	stab := b.scopeDict.Data[parent.ID]
	fmt.Println("{", stab, "}")

	fmt.Println(instr.Local)
}

func (b *borrowChecker) validate(fn *ir.Function) {
	b.block = fn.Body
	fmt.Println("pushed main")
	b.visitBlock(b.block)
}

func borrowCheck(mod *ir.Module, scopeDict *ir.ScopeDict) []api.CompilerError {
	errs := []api.CompilerError{}

	for name, fn := range mod.Functions {
		checker := &borrowChecker{
			scopeDict: scopeDict,
		}
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

	var scopeDict ir.ScopeDict
	if err := jsoniter.Unmarshal([]byte(req.ScopeMap), &scopeDict); err != nil {
		panic(err)
	}

	errs := borrowCheck(&irMod, &scopeDict)

	resp := api.KrugResponse{
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}
