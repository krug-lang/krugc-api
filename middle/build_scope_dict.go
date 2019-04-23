package middle

import (
	"fmt"

	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
)

type scopeDictBuilder struct {
	mod        *ir.Module
	curr       *ir.SymbolTable
	outer      []*ir.SymbolTable
	errs       []api.CompilerError
	blockCount int
	scopeDict  *ir.ScopeDict
}

func (b *scopeDictBuilder) error(err api.CompilerError) {
	b.errs = append(b.errs, err)
}

func (b *scopeDictBuilder) pushStab(name string) *ir.SymbolTable {
	// store the previous stab
	prev := b.curr

	// setup a new stab to push
	pushed := ir.NewSymbolTable(prev)
	// store as the current
	b.curr = pushed

	// append this stab to the previous
	// inner blocks.
	if prev != nil {
		prev.Inner = append(prev.Inner, pushed)
	}

	// set the outer to the previous
	b.outer = append(b.outer, prev)

	return pushed
}

func (b *scopeDictBuilder) popStab() *ir.SymbolTable {
	// store the stab to pop
	popped := b.curr

	// if we dont have a previous scope
	// to pop to, nullify the current scope
	// and return
	if len(b.outer) == 0 {
		b.curr = nil
		return popped
	}

	// set the current scope to be the old scope
	b.curr = b.outer[len(b.outer)-1]
	b.outer = b.outer[:len(b.outer)-1]

	return popped
}

func (b *scopeDictBuilder) pushBlock(id int) *ir.SymbolTable {
	return b.pushStab(fmt.Sprintf("%d", id))
}

// TODO generate some kind of key that maps an ir.Block
// to the symbol table.
func (b *scopeDictBuilder) visitBlock(i *ir.Block) *ir.SymbolTable {
	i.Stab = b.pushBlock(b.blockCount)
	b.assign(i, i.Stab)
	b.blockCount++

	for _, instr := range i.Instr {
		b.visitInstr(instr)
	}

	b.popStab()
	return i.Stab
}

func (b *scopeDictBuilder) visitIfStat(iff *ir.IfStatement) {
	b.visitBlock(iff.True)
	for _, e := range iff.ElseIf {
		b.visitBlock(e.Body)
	}
	if iff.Else != nil {
		b.visitBlock(iff.Else)
	}
}

func (b *scopeDictBuilder) visitWhileLoop(while *ir.WhileLoop) {
	b.visitBlock(while.Body)
}

func (b *scopeDictBuilder) visitLoop(loop *ir.Loop) {
	b.visitBlock(loop.Body)
}

func (b *scopeDictBuilder) visitInstr(i *ir.Instruction) {
	switch i.Kind {

	// TODO(FElix):
	// store owned in the alloc, and local instr

	case ir.AllocaInstr:
		instr := i.Alloca
		ok := b.curr.Register(instr.Name.Value, &ir.SymbolValue{
			Kind:   ir.SymbolKind,
			Symbol: ir.NewSymbol(instr.Name, instr.Owned, instr.Mutable),
		})
		if !ok {
			b.error(api.NewSymbolError(instr.Name.Value, instr.Name.Span...))
		}

	case ir.LocalInstr:
		instr := i.Local
		ok := b.curr.Register(instr.Name.Value, &ir.SymbolValue{
			Kind:   ir.SymbolKind,
			Symbol: ir.NewSymbol(instr.Name, instr.Owned, instr.Mutable),
		})
		if !ok {
			b.error(api.NewSymbolError(instr.Name.Value, instr.Name.Span...))
		}

	case ir.IfStatementInstr:
		instr := i.IfStatement
		b.visitIfStat(instr)
	case ir.WhileLoopInstr:
		instr := i.WhileLoop
		b.visitWhileLoop(instr)
	case ir.LoopInstr:
		instr := i.Loop
		b.visitLoop(instr)

	case ir.BlockInstr:
		instr := i.Block
		b.visitBlock(instr)

	case ir.LabelInstr:
	case ir.JumpInstr:
	case ir.ReturnInstr:
		// nop

	case ir.ExpressionInstr:
		fmt.Println(i.ExpressionStatement)

	/*
		case *ir.Path:
			return
		case *ir.Return:
			return
		case *ir.Call:
			return
		case *ir.Assign:
			return
	*/

	default:
		panic(fmt.Sprintf("unhandled instr %s", i.Kind))
	}
}

func (b *scopeDictBuilder) assign(block *ir.Block, stab *ir.SymbolTable) {
	b.scopeDict.Data[block.ID] = stab
}

// hack we shouldnt have to do this in the first place?
func (b *scopeDictBuilder) clearScope() {
	b.curr = nil
}

func (b *scopeDictBuilder) visitFunc(fn *ir.Function) *ir.SymbolTable {
	b.clearScope()
	res := b.pushStab(fn.Name.Value)
	b.assign(fn.Body, res)

	// reset the block count
	b.blockCount = 0

	// introduce params into the function scope.
	for _, name := range fn.Param.Order {
		param := fn.Param.Data[name.Value]
		ok := b.curr.Register(name.Value, &ir.SymbolValue{
			Kind:   ir.SymbolKind,
			Symbol: ir.NewSymbol(name, param.Owned, param.Mutable),
		})
		if !ok {
			b.error(api.NewSymbolError(name.Value, name.Span...))
		}
	}

	// manually visit the functions body as we've
	// already pushed a scope.
	for _, instr := range fn.Body.Instr {
		b.visitInstr(instr)
	}

	b.popStab()

	return res
}

func (b *scopeDictBuilder) visitStructure(s *ir.Structure) *ir.SymbolTable {
	stab := b.pushStab(s.Name.Value)

	for _, name := range s.Fields.Order {
		ok := stab.Register(name.Value, &ir.SymbolValue{
			Kind: ir.SymbolKind,
			// structure fields are:
			// - not owners of their memory (? FIXME)
			// - mutable
			Symbol: ir.NewSymbol(name, false, true),
		})
		if !ok {
			b.error(api.NewSymbolError(name.Value, name.Span...))
		}
	}

	return stab
}

func buildScopeDict(mod *ir.Module) (*ir.ScopeDict, []api.CompilerError) {
	b := &scopeDictBuilder{
		mod,
		nil,
		[]*ir.SymbolTable{},
		[]api.CompilerError{},
		0,
		ir.NewScopeDict(),
	}

	for _, fn := range mod.Functions {
		b.visitFunc(fn)
	}

	return b.scopeDict, b.errs
}
