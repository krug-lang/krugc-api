package middle

import (
	"fmt"
	"reflect"

	"github.com/krug-lang/server/api"
	"github.com/krug-lang/server/ir"
)

type builder struct {
	mod        *ir.Module
	curr       *ir.SymbolTable
	errs       []api.CompilerError
	blockCount int
}

func (b *builder) error(err api.CompilerError) {
	b.errs = append(b.errs, err)
}

func (b *builder) pushStab(name string) *ir.SymbolTable {
	old := b.curr
	b.curr = ir.NewSymbolTable(old)
	return b.curr
}

func (b *builder) popStab() *ir.SymbolTable {
	current := b.curr
	if current.Outer == nil {
		panic("no stab to pop to")
	}
	b.curr = current.Outer
	return current
}

func (b *builder) pushBlock(id int) *ir.SymbolTable {
	return b.pushStab(fmt.Sprintf("%s", id))
}

func (b *builder) visitBlock(i *ir.Block) *ir.SymbolTable {
	i.Stab = b.pushBlock(b.blockCount)
	b.blockCount++

	for _, instr := range i.Instr {
		b.visitInstr(instr)
	}

	b.popStab()
	return i.Stab
}

func (b *builder) visitIfStat(iff *ir.IfStatement) {
	b.visitBlock(iff.True)
	for _, e := range iff.ElseIf {
		b.visitBlock(e.Body)
	}
	if iff.Else != nil {
		b.visitBlock(iff.Else)
	}
}

func (b *builder) visitWhileLoop(while *ir.WhileLoop) {
	b.visitBlock(while.Body)
}

func (b *builder) visitLoop(loop *ir.Loop) {
	b.visitBlock(loop.Body)
}

func (b *builder) visitInstr(i ir.Instruction) {
	switch instr := i.(type) {
	case *ir.Alloca:
		ok := b.curr.Register(instr.Name.Value, ir.NewSymbol(instr.Name))
		if !ok {
			b.error(api.NewSymbolError(instr.Name.Value, instr.Name.Span...))
		}

	case *ir.Local:
		ok := b.curr.Register(instr.Name.Value, ir.NewSymbol(instr.Name))
		if !ok {
			b.error(api.NewSymbolError(instr.Name.Value, instr.Name.Span...))
		}

	case *ir.IfStatement:
		b.visitIfStat(instr)
	case *ir.WhileLoop:
		b.visitWhileLoop(instr)
	case *ir.Loop:
		b.visitLoop(instr)

	case *ir.Block:
		b.visitBlock(instr)

	case *ir.Path:
		return
	case *ir.Return:
		return
	case *ir.Call:
		return
	case *ir.Assign:
		return

	default:
		panic(fmt.Sprintf("unhandled instr %s", reflect.TypeOf(instr)))
	}
}

func (b *builder) visitFunc(fn *ir.Function) *ir.SymbolTable {
	res := b.pushStab(fn.Name.Value)

	// reset the block count
	b.blockCount = 0

	// introduce params into the function scope.
	for _, name := range fn.Param.Order {
		ok := b.curr.Register(name.Value, ir.NewSymbol(name))
		if !ok {
			b.error(api.NewSymbolError(name.Value, name.Span...))
		}
	}

	// manually visit the functions body as we've
	// already pushed a scope.
	for _, instr := range fn.Body.Instr {
		b.visitInstr(instr)
	}

	return res
}

func (b *builder) visitStructure(s *ir.Structure) *ir.SymbolTable {
	stab := b.pushStab(s.Name.Value)

	for _, name := range s.Fields.Order {
		ok := stab.Register(name.Value, ir.NewSymbol(name))
		if !ok {
			b.error(api.NewSymbolError(name.Value, name.Span...))
		}
	}

	return stab
}

func buildScope(mod *ir.Module) (*ir.ScopeMap, []api.CompilerError) {
	b := &builder{
		mod,
		nil,
		[]api.CompilerError{},
		0,
	}

	scopeMap := ir.NewScopeMap()

	// create stabs for the structs
	for _, structure := range mod.Structures {
		stab := b.visitStructure(structure)

		ok := scopeMap.RegisterStructure(structure.Name.Value, stab)
		if !ok {
			b.error(api.NewSymbolError(structure.Name.Value, structure.Name.Span...))
		}
	}

	// create stabs for the functions
	for _, fn := range mod.Functions {
		stab := b.visitFunc(fn)

		ok := scopeMap.RegisterFunction(fn.Name.Value, stab)
		if !ok {
			b.error(api.NewSymbolError(fn.Name.Value, fn.Name.Span...))
		}
	}

	/*
		// create stabs for the impls
		for name, impl := range mod.Impls {
			impl.Stab = b.pushStab(name)

			for _, fn := range impl.Methods {
				fn.Stab = b.visitFunc(fn)
			}

			b.popStab()
		}
	*/

	return scopeMap, b.errs
}
