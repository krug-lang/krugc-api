package middle

import (
	"fmt"
	"reflect"

	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
)

type symResolvePass struct {
	mod    *ir.Module
	errors []api.CompilerError
	curr   *ir.SymbolTable
}

func (s *symResolvePass) error(err api.CompilerError) {
	s.errors = append(s.errors, err)
}

func (s *symResolvePass) push(stab *ir.SymbolTable) {
	s.curr = stab
}

func (s *symResolvePass) pop() {
	if s.curr != nil {
		s.curr = s.curr.Outer
	}
}

func (s *symResolvePass) resolveIden(i *ir.Identifier) (ir.SymbolValue, bool) {
	val, ok := s.curr.Lookup(i.Name.Value)
	if !ok {
		s.error(api.NewUnresolvedSymbol(i.Name.Value, i.Name.Span...))
	}
	return val, ok
}

func (s *symResolvePass) resolveAssign(a *ir.Assign) {
	s.resolveValue(a.LHand)
	s.resolveValue(a.RHand)
}

func (s *symResolvePass) resolveValue(e *ir.Value) ir.SymbolValue {
	switch e.Kind {
	case ir.IntegerValueValue:
		return nil
	case ir.StringValueValue:
		return nil
	case ir.FloatingValueValue:
		return nil

	case ir.AssignValue:
		s.resolveAssign(e.Assign)
		return nil

	case ir.BinaryExpressionValue:
		s.resolveValue(e.BinaryExpression.LHand)
		s.resolveValue(e.BinaryExpression.RHand)
		return nil

	case ir.IdentifierValue:
		stab, _ := s.resolveIden(e.Identifier)
		return stab

	default:
		panic(fmt.Sprintf("unhandled val %s", reflect.TypeOf(e)))
	}
}

func (s *symResolvePass) resolveAlloca(v *ir.Alloca) {
	if v.Val != nil {
		s.resolveValue(v.Val)
	}
}

func (s *symResolvePass) resolveLocal(v *ir.Local) {
	if v.Val != nil {
		s.resolveValue(v.Val)
	}
}

func (s *symResolvePass) resolveBlock(b *ir.Block) {
	s.push(b.Stab)
	for _, instr := range b.Instr {
		s.resolveInstr(instr)
	}
	s.pop()
}

func (s *symResolvePass) resolveCall(c *ir.Call) {
	// TODO:
}

func (s *symResolvePass) resolveValueVia(last *ir.SymbolTable, val *ir.Value) ir.SymbolValue {
	if last == nil {
		return s.resolveValue(val)
	}

	switch val.Kind {
	case ir.IdentifierValue:
		v := val.Identifier
		val, ok := last.Lookup(v.Name.Value)
		if !ok {
			s.error(api.NewUnresolvedSymbol(v.Name.Value))
		}
		return val

	default:
		panic(fmt.Sprintf("unhandled value %s", reflect.TypeOf(val)))
	}
}

func (s *symResolvePass) resolvePath(p *ir.Path) {
	s.resolveValue(p.Values[0])
}

func (s *symResolvePass) resolveInstr(i *ir.Instruction) {
	switch i.Kind {
	case ir.AllocaInstr:
		instr := i.Alloca
		s.resolveAlloca(instr)
	case ir.LocalInstr:
		instr := i.Local
		s.resolveLocal(instr)

	case ir.ReturnInstr:
		instr := i.Return
		if instr.Val != nil {
			s.resolveValue(instr.Val)
		}

	// FIXME
	case ir.PathValue:
		// s.resolvePath(instr)

	case ir.IfStatementInstr:
		instr := i.IfStatement
		s.resolveValue(instr.Cond)
		s.resolveBlock(instr.True)
		for _, e := range instr.ElseIf {
			s.resolveBlock(e.Body)
		}
		if instr.Else != nil {
			s.resolveBlock(instr.Else)
		}
	case ir.WhileLoopInstr:
		instr := i.WhileLoop
		s.resolveValue(instr.Cond)
		if instr.Post != nil {
			s.resolveValue(instr.Post)
		}
		s.resolveBlock(instr.Body)
	case ir.LoopInstr:
		instr := i.Loop
		s.resolveBlock(instr.Body)

	case ir.BlockInstr:
		instr := i.Block
		s.resolveBlock(instr)

	// FIXME
	case ir.CallValue:
		// s.resolveCall(instr)

	case ir.AssignInstr:
		instr := i.Assign
		s.resolveAssign(instr)

	default:
		panic(fmt.Sprintf("unhandled instr %s", reflect.TypeOf(i)))
	}
}

func (s *symResolvePass) resolveFunc(fn *ir.Function) {
	s.push(fn.Stab)

	/*
			json, err := jsoniter.MarshalIndent(s.curr, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(json))
	*/

	for _, instr := range fn.Body.Instr {
		s.resolveInstr(instr)
	}

	s.pop()
}

func symResolve(mod *ir.Module) (*ir.Module, []api.CompilerError) {
	srp := &symResolvePass{mod, []api.CompilerError{}, nil}

	for _, impl := range mod.Impls {
		for _, method := range impl.Methods {
			srp.resolveFunc(method)
		}
	}

	for _, fn := range mod.Functions {
		srp.resolveFunc(fn)
	}

	return srp.mod, srp.errors
}
