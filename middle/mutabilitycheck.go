package middle

import (
	"fmt"
	"github.com/krug-lang/caasper/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/caasper/api"
	"github.com/krug-lang/caasper/ir"
	jsoniter "github.com/json-iterator/go"
)

type mutChecker struct {
	mod  *ir.Module
	dict *ir.ScopeDict
	errs []api.CompilerError
}

func (m *mutChecker) error(e api.CompilerError) {
	m.errs = append(m.errs, e)
}

func isSymbolMutable(sym *ir.SymbolValue) bool {
	switch sym.Kind {
	case ir.SymbolKind:
		return sym.Symbol.Mutable
	}
	return false
}

func (m *mutChecker) checkIdenMutable(parent *ir.Block, iden *ir.Identifier) bool {
	dict, _ := m.dict.Data[parent.ID]

	if sym, ok := dict.Lookup(iden.Name.Value); ok {
		return isSymbolMutable(sym)
	}

	// check parent sym table?
	return false
}

func (m *mutChecker) checkMutable(parent *ir.Block, val *ir.Value) bool {
	switch val.Kind {
	case ir.IdentifierValue:
		mutable := m.checkIdenMutable(parent, val.Identifier)
		if !mutable {
			tok := val.Identifier.Name
			m.error(api.NewMutabilityError(tok.Value, tok.Span...))
		}
		return mutable

	case ir.IntegerValueValue:
	case ir.StringValueValue:
	case ir.CharacterValueValue:
	case ir.FloatingValueValue:

	default:
		panic(fmt.Sprintf("checkMutable: unhandled value %s", val.Kind))
	}

	return false
}

func (m *mutChecker) visitIden(parent *ir.Block, iden *ir.Identifier) {
}

func (m *mutChecker) visitAssign(parent *ir.Block, assign *ir.Assign) {
	m.checkMutable(parent, assign.LHand)
}

func (m *mutChecker) visitExpr(parent *ir.Block, expr *ir.Value) {
	switch expr.Kind {
	case ir.IdentifierValue:
		m.visitIden(parent, expr.Identifier)

	case ir.UnaryExpressionValue:
		m.visitExpr(parent, expr.UnaryExpression.Val)

	case ir.BinaryExpressionValue:
		m.visitExpr(parent, expr.BinaryExpression.LHand)
		m.visitExpr(parent, expr.BinaryExpression.RHand)

	case ir.AssignValue:
		m.visitAssign(parent, expr.Assign)

	case ir.CallValue:
		// TODO
		// we need to ensure that the call sites
		// paramteres share the same mutability as the
		// values passed through?

	case ir.IntegerValueValue:
	case ir.StringValueValue:
	case ir.CharacterValueValue:
	case ir.FloatingValueValue:

	case ir.PathValue:
		// TODO!

	default:
		panic(fmt.Sprintf("unhandled expr %s", expr.Kind))
	}
}

func (m *mutChecker) visitLocal(parent *ir.Block, local *ir.Local) {
	// nop?
}

func (m *mutChecker) visitInstr(parent *ir.Block, instr *ir.Instruction) {
	switch instr.Kind {
	case ir.LocalInstr:
		m.visitLocal(parent, instr.Local)

	case ir.ExpressionInstr:
		m.visitExpr(parent, instr.ExpressionStatement)

	case ir.ReturnInstr:
		if val := instr.Return.Val; val != nil {
			m.visitExpr(parent, val)
		}

	case ir.IfStatementInstr:
		iff := instr.IfStatement
		m.visitExpr(parent, iff.Cond)
		m.visitBlock(iff.True)
		for _, elif := range iff.ElseIf {
			m.visitBlock(elif.Body)
			// FIXME hm what block is this a part of?
			m.visitExpr(elif.Body, elif.Cond)
		}
		if eb := iff.Else; eb != nil {
			m.visitBlock(eb)
		}

	case ir.LoopInstr:
		m.visitBlock(instr.Loop.Body)

	case ir.WhileLoopInstr:
		wl := instr.WhileLoop
		m.visitExpr(wl.Body, wl.Cond)
		if post := wl.Post; post != nil {
			m.visitExpr(wl.Body, post)
		}
		m.visitBlock(wl.Body)

	case ir.LabelInstr:
	case ir.JumpInstr:
		// nop

	default:
		panic(fmt.Sprintf("unhandled instr %s", instr.Kind))
	}
}

func (m *mutChecker) visitBlock(block *ir.Block) {
	for _, instr := range block.Instr {
		switch instr.Kind {
		case ir.BlockInstr:
			m.visitBlock(instr.Block)
		default:
			m.visitInstr(block, instr)
		}
	}
}

func (m *mutChecker) check(fn *ir.Function) {
	m.visitBlock(fn.Body)
}

func mutCheck(mod *ir.Module, dict *ir.ScopeDict) []api.CompilerError {
	checker := &mutChecker{
		mod, dict, []api.CompilerError{},
	}

	for _, name := range mod.FunctionOrder {
		fn, _ := mod.Functions[name.Value]
		checker.check(fn)
	}

	return checker.errs
}

func MutabilityCheck(c *gin.Context) {
	var req entity.MutabilityCheckRequest
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

	errs := mutCheck(&irMod, &scopeDict)

	resp := entity.KrugResponse{
		Data:   "",
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}
