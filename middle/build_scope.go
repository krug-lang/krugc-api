package middle

import (
	"fmt"
	"math/rand"

	"github.com/krug-lang/krugc-api/api"
	"github.com/krug-lang/krugc-api/ir"
)

type SymbolValue interface{}

type Symbol struct {
	Name string
}

type SymbolTable struct {
	Id      int
	Outer   *SymbolTable
	Symbols map[string]SymbolValue
}

// hm
func (s *SymbolTable) Register(name string, sym SymbolValue) {
	s.Symbols[name] = sym
}

func (s *SymbolTable) Lookup(name string) (SymbolValue, bool) {
	if sym, ok := s.Symbols[name]; ok {
		if ok {
			return sym, ok
		}
	}

	if s.Outer == nil {
		return nil, false
	}

	return s.Outer.Lookup(name)
}

func newSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Id:      rand.Intn(30000),
		Outer:   outer,
		Symbols: map[string]SymbolValue{},
	}
}

type builder struct {
	mod  *ir.Module
	curr *SymbolTable
	errs []api.CompilerError
}

func (b *builder) pushStab(name string) *SymbolTable {
	old := b.curr
	b.curr = newSymbolTable(old)
	return b.curr
}

func (b *builder) popStab() *SymbolTable {
	current := b.curr
	if current.Outer == nil {
		panic("no stab to pop to")
	}
	b.curr = current.Outer
	return current
}

func (b *builder) pushBlock() *SymbolTable {
	id := 0
	return b.pushStab(fmt.Sprintf("%s", id))
}

func buildScope(mod *ir.Module) (*SymbolTable, []api.CompilerError) {
	b := &builder{
		mod,
		nil,
		[]api.CompilerError{},
	}

	root := b.pushStab("0_global")

	return root, b.errs
}
