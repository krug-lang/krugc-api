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
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/api"
	"github.com/krug-lang/caasper/ir"
)

type lifetime struct {
	outer  *lifetime
	id     int
	locals map[string]*local
}

func (l *lifetime) findLocal(name string) *local {
	// find in this lifetime
	if loc, ok := l.locals[name]; ok {
		return loc
	}

	// no outer lifetime? couldnt find it
	if l.outer == nil {
		return nil
	}

	// look in outer lifetime
	return l.outer.findLocal(name)
}

func (l *lifetime) addLocal(loc *local) {
	fmt.Println("addLocal ", loc.loc.Name.Value)
	// LOL
	l.locals[loc.loc.Name.Value] = loc
}

func newLifetime(outer *lifetime) *lifetime {
	return &lifetime{
		outer:  outer,
		id:     0,
		locals: map[string]*local{},
	}
}

type local struct {
	loc       *ir.Local
	loans     []*ir.Value
	borrowers []*ir.Value
}

func newLoc(loc *ir.Local) *local {
	return &local{loc, []*ir.Value{}, []*ir.Value{}}
}

func (l *local) loanTo(v *ir.Value) {
	fmt.Println("loaning", l.loc.Name.Value)
	l.loans = append(l.loans, v)
}

func (l *local) addBorrower(v *ir.Value) {
	fmt.Println("borrow", l.loc.Name.Value)
	l.borrowers = append(l.borrowers, v)
}

type borrowChecker struct {
	errs          []api.CompilerError
	block         *ir.Block
	scopeDict     *ir.ScopeDict
	lifetime      *lifetime
	prevLifetimes []*lifetime
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

func (b *borrowChecker) error(err api.CompilerError) {
	b.errs = append(b.errs, err)
}

func (b *borrowChecker) pushLifetime() {
	newLifetime := newLifetime(b.lifetime)
	b.lifetime = newLifetime
}

func (b *borrowChecker) visitLocal(loc *ir.Local) {
	if loc.Val != nil {
		b.visitExpr(loc.Val, loc.Val)
	}

	// this variable owns its memory,
	// add it to the lifetime.
	if loc.Owned {
		lt := b.lifetime
		lt.addLocal(newLoc(loc))
	}
}

func (b *borrowChecker) visitCall(call *ir.Call) {
	// we dont care about the call sites type information
	// we can check for params that are identifiers
	// and see if they are using owned values or not
	for _, param := range call.Params {
		b.visitExpr(call.Left, param)
	}
}

// getIdentifierRef will look for the local that this
// identifier is referencing in the lifetimes or parent lifetimes
func (b *borrowChecker) getIdentifierRef(iden *ir.Identifier) *local {
	owner := b.lifetime.findLocal(iden.Name.Value)
	if owner != nil {
		return owner
	}
	return nil
}

func (b *borrowChecker) visitBuiltin(lhand *ir.Value, builtin *ir.Builtin) {
	if strings.Compare(builtin.Name, "ref") != 0 {
		return
	}

	owner := b.getIdentifierRef(builtin.Iden)
	if owner != nil && lhand != nil {
		// as many borrowers as we want.
		// these borrows are _immutable_
		owner.addBorrower(lhand)
	}
}

func (b *borrowChecker) tryLoan(iden *ir.Identifier, loanee *ir.Value) {
	owner := b.getIdentifierRef(iden)
	if owner != nil {

		// TODO this error could also show who the value has been
		// loaned to prior?
		if len(owner.loans) >= 1 {
			b.error(api.NewMovedValueError(iden.Name.Value, iden.Name.Span...))
		}
	}

	// sometimes the loanee can be nil as we
	// dont know who the loan is going to, e.g.
	// return X
	// we loan X to the function return val?
	if owner != nil && loanee != nil {
		owner.loanTo(loanee)
	}
}

func (b *borrowChecker) visitExpr(lhand *ir.Value, expr *ir.Value) {
	switch expr.Kind {
	case ir.CallValue:
		b.visitCall(expr.Call)
	case ir.IdentifierValue:
		b.tryLoan(expr.Identifier, lhand)

	case ir.UnaryExpressionValue:
		b.visitExpr(lhand, expr.UnaryExpression.Val)

	case ir.BinaryExpressionValue:
		b.visitExpr(lhand, expr.BinaryExpression.LHand)
		b.visitExpr(lhand, expr.BinaryExpression.RHand)

	case ir.IntegerValueValue:
		break
	case ir.FloatingValueValue:
		break
	case ir.StringValueValue:
		break

	case ir.BuiltinValue:
		b.visitBuiltin(lhand, expr.Builtin)

	case ir.AssignValue:
		b.visitExpr(expr.Assign.LHand, expr.Assign.RHand)

	case ir.PathValue:
		// TODO!

	default:
		panic(fmt.Sprintf("unhandled expr %s", expr.Kind))
	}
}

func (b *borrowChecker) visitInstr(parent *ir.Block, instr *ir.Instruction) {
	stab := b.scopeDict.Data[parent.ID]
	fmt.Println("{", stab, "}")

	switch instr.Kind {
	case ir.LocalInstr:
		b.visitLocal(instr.Local)
	case ir.ExpressionInstr:
		b.visitExpr(nil, instr.ExpressionStatement)
	case ir.ReturnInstr:
		ret := instr.Return
		if ret.Val != nil {
			b.visitExpr(nil, ret.Val)
		}

	case ir.WhileLoopInstr:
		b.visitExpr(nil, instr.WhileLoop.Cond)
		if post := instr.WhileLoop.Post; post != nil {
			b.visitExpr(nil, post)
		}
		b.visitBlock(instr.WhileLoop.Body)

	case ir.LoopInstr:
		b.visitBlock(instr.Loop.Body)

	case ir.BreakInstr:
		break

	case ir.IfStatementInstr:
		b.visitExpr(nil, instr.IfStatement.Cond)
		b.visitBlock(instr.IfStatement.True)
		// TODO else if and elses.

	case ir.JumpInstr:
	case ir.LabelInstr:

	default:
		panic(fmt.Sprintf("unhandled instruction %s", instr.Kind))
	}
}

func (b *borrowChecker) popLifetime() {
	b.lifetime = b.lifetime.outer
}

func (b *borrowChecker) visitBlockPre(block *ir.Block, preVisit func()) {
	b.pushLifetime()

	if preVisit != nil {
		preVisit()
	}

	for _, instr := range block.Instr {
		switch instr.Kind {
		case ir.BlockInstr:
			b.visitBlock(instr.Block)
		default:
			b.visitInstr(block, instr)
		}
	}

	b.popLifetime()
}

func (b *borrowChecker) visitBlock(block *ir.Block) {
	b.visitBlockPre(block, nil)
}

func (b *borrowChecker) validate(fn *ir.Function) {
	b.block = fn.Body

	b.visitBlockPre(b.block, func() {
		for _, tok := range fn.Param.Order {
			loc, ok := fn.Param.Data[tok.Value]
			if !ok {
				panic("this should never happen")
			}
			if loc.Owned {
				b.lifetime.addLocal(newLoc(loc))
			}
		}
	})
}

func borrowCheck(mod *ir.Module, scopeDict *ir.ScopeDict) []api.CompilerError {
	errs := []api.CompilerError{}

	for name, fn := range mod.Functions {
		checker := &borrowChecker{
			scopeDict: scopeDict,
		}
		fmt.Println("validating ", name)
		checker.validate(fn)

		errs = append(errs, checker.errs...)
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
		Data:   "",
		Errors: errs,
	}

	c.JSON(http.StatusOK, &resp)
}
