package middle

import (
	"fmt"
	"reflect"

	"github.com/krug-lang/krugc-api/api"
	"github.com/krug-lang/krugc-api/ir"
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

func (s *symResolvePass) resolveValue(e ir.Value) ir.SymbolValue {
	switch expr := e.(type) {
	case *ir.IntegerValue:
		return nil
	case *ir.StringValue:
		return nil
	case *ir.FloatingValue:
		return nil

	case *ir.Assign:
		s.resolveAssign(expr)
		return nil

	case *ir.BinaryExpression:
		s.resolveValue(expr.LHand)
		s.resolveValue(expr.RHand)
		return nil

	case *ir.Identifier:
		stab, _ := s.resolveIden(expr)
		return stab

	default:
		panic(fmt.Sprintf("unhandled val %s", reflect.TypeOf(expr)))
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

func (s *symResolvePass) resolveValueVia(last *ir.SymbolTable, val ir.Value) ir.SymbolValue {
	if last == nil {
		return s.resolveValue(val)
	}

	switch v := val.(type) {
	case *ir.Identifier:
		val, ok := last.Lookup(v.Name.Value)
		if !ok {
			s.error(api.NewUnresolvedSymbol(v.Name.Value))
		}
		return val

	default:
		panic(fmt.Sprintf("unhandled value %s", reflect.TypeOf(v)))
	}
}

func (s *symResolvePass) resolvePath(p *ir.Path) {
	s.resolveValue(p.Values[0])
}

func (s *symResolvePass) resolveInstr(i ir.Instruction) {
	switch instr := i.(type) {
	case *ir.Alloca:
		s.resolveAlloca(instr)
	case *ir.Local:
		s.resolveLocal(instr)

	case *ir.Return:
		if instr.Val != nil {
			s.resolveValue(instr.Val)
		}
	case *ir.Path:
		s.resolvePath(instr)

	case *ir.IfStatement:
		s.resolveValue(instr.Cond)
		s.resolveBlock(instr.True)
		for _, e := range instr.ElseIf {
			s.resolveBlock(e.Body)
		}
		if instr.Else != nil {
			s.resolveBlock(instr.Else)
		}
	case *ir.WhileLoop:
		s.resolveValue(instr.Cond)
		if instr.Post != nil {
			s.resolveValue(instr.Post)
		}
		s.resolveBlock(instr.Body)
	case *ir.Loop:
		s.resolveBlock(instr.Body)

	case *ir.Block:
		s.resolveBlock(instr)

	case *ir.Call:
		s.resolveCall(instr)
	case *ir.Assign:
		s.resolveAssign(instr)

	default:
		panic(fmt.Sprintf("unhandled instr %s", reflect.TypeOf(instr)))
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
