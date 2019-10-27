package middle

import (
	"fmt"
	"strings"

	"github.com/krug-lang/caasper/api"
	"github.com/krug-lang/caasper/ir"
)

/*
	this pass will build a control flow graph of the module
	and check for any unused functions.
*/

type functionGraph struct {
	nodes   []*ir.Function
	nodeSet map[string]*ir.Function
	edges   map[string][]*ir.Function
}

func (f *functionGraph) hasNode(name string) bool {
	_, ok := f.nodeSet[name]
	return ok
}

func (f *functionGraph) addNode(nodes ...*ir.Function) {
	if f.nodeSet == nil {
		f.nodeSet = map[string]*ir.Function{}
	}

	for _, node := range nodes {
		f.nodes = append(f.nodes, node)
		f.nodeSet[node.Name.Value] = node
	}
}

// add an edge from <- to
// i.e. the to function has an edge of where it is called from
func (f *functionGraph) addEdge(from, to string) {
	if f.edges == nil {
		f.edges = make(map[string][]*ir.Function)
	}
	fromFn := f.nodeSet[from]
	f.edges[to] = append(f.edges[to], fromFn)
}

type visitor struct {
	g        *functionGraph
	mod      *ir.Module
	errs     []api.CompilerError
	currFunc *ir.Function
}

func (v *visitor) err(err api.CompilerError) {
	v.errs = append(v.errs, err)
}

func (v *visitor) resolveCall(call *ir.Call) {
	if iden := call.Left.Identifier; iden != nil {
		fn := iden.Name.Value
		if v.g.hasNode(fn) {
			v.g.addEdge(v.currFunc.Name.Value, fn)
			fmt.Println("We can resolve the function ", fn)
		}
	}
}

func (v *visitor) visitValue(expr *ir.Value) {
	switch expr.Kind {
	case ir.CallValue:
		v.resolveCall(expr.Call)

	case ir.BinaryExpressionValue:
		v.visitValue(expr.BinaryExpression.LHand)
		v.visitValue(expr.BinaryExpression.RHand)

	case ir.FloatingValueValue:
		fallthrough
	case ir.StringValueValue:
		fallthrough
	case ir.CharacterValueValue:
		fallthrough
	case ir.IdentifierValue:
		fallthrough
	case ir.IntegerValueValue:
		return

	default:
		v.err(api.NewUnimplementedError("unused_func", fmt.Sprintf("visitValue:%s", expr.Kind)))
	}
}

func (v *visitor) visitLocal(local *ir.Local) {
	if val := local.Val; val != nil {
		v.visitValue(val)
	}
}

func (v *visitor) visitIf(instr *ir.IfStatement) {
	v.visitValue(instr.Cond)

	for _, elif := range instr.ElseIf {
		v.visitValue(elif.Cond)
		v.visitBlock(elif.Body)
	}

	if els := instr.Else; els != nil {
		v.visitBlock(els)
	}
}

func (v *visitor) visitInstr(instr *ir.Instruction) {
	switch instr.Kind {
	case ir.ExpressionInstr:
		v.visitValue(instr.ExpressionStatement)
	case ir.LocalInstr:
		v.visitLocal(instr.Local)

	case ir.ReturnInstr:
		if ret := instr.Return.Val; ret != nil {
			v.visitValue(ret)
		}
	case ir.BlockInstr:
		v.visitBlock(instr.Block)
	case ir.LoopInstr:
		v.visitBlock(instr.Loop.Body)
	case ir.IfStatementInstr:
		v.visitIf(instr.IfStatement)

	default:
		v.err(api.NewUnimplementedError("unused_func", fmt.Sprintf("visitInstr:%s", instr.Kind)))
	}
}

func (v *visitor) visitBlock(b *ir.Block) {
	for _, stat := range b.Instr {
		v.visitInstr(stat)
	}
	for _, stat := range b.DeferStack {
		if stat.Block != nil {
			v.visitBlock(stat.Block)
		} else if stat.Stat != nil {
			v.visitInstr(stat.Stat)
		}
	}
}

func calculateFunctionUsage(g *functionGraph, mod *ir.Module) []api.CompilerError {
	v := visitor{g, mod, []api.CompilerError{}, nil}

	for _, fn := range mod.Functions {
		v.currFunc = fn
		v.visitBlock(fn.Body)
	}

	return v.errs
}

func UnusedFunc(mod *ir.Module, scope *ir.ScopeDict) []api.CompilerError {
	g := &functionGraph{}
	for _, fn := range mod.Functions {
		g.addNode(fn)
	}

	errs := calculateFunctionUsage(g, mod)

	for _, fn := range mod.Functions {
		fnName := fn.Name.Value

		// main will be unused as its invoked
		// from outside the language
		if strings.Compare(fnName, "main") == 0 {
			continue
		}

		edges := g.edges[fnName]
		if len(edges) == 0 {
			errs = append(errs, api.NewUnusedFunction(fn.Name.Value, fn.Name.Span...))
		}
	}

	return errs
}