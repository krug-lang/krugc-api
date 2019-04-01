package ir

import (
	"encoding/gob"
	"math/rand"

	"github.com/hugobrains/krug-serv/front"
)

func init() {
	gob.Register(&SymbolTable{})
	gob.Register(&Symbol{})
}

type SymbolValue interface {
	SymbolTypeName() string
	GetType() Type
}

type Symbol struct {
	Name front.Token
}

func (s *Symbol) GetType() Type {
	return nil
}

func (s *Symbol) SymbolTypeName() string {
	return "symbol"
}

func NewSymbol(name front.Token) *Symbol {
	return &Symbol{name}
}

type SymbolTable struct {
	Id      int
	Outer   *SymbolTable
	Types   map[string]Type
	Symbols map[string]SymbolValue
}

func (s *SymbolTable) String() string {
	res := "{"

	idx := 0
	for _, sym := range s.Symbols {
		if idx != 0 {
			res += " "
		}

		switch sy := sym.(type) {
		case *Symbol:
			res += sy.Name.Value
		}

		idx++
	}
	res += "}"
	return res
}

func (s *SymbolTable) RegisterType(name string, t Type) {
	s.Types[name] = t
}

func (s *SymbolTable) LookupType(name string) (Type, bool) {
	typ, ok := s.Types[name]
	// TODO look in outer scope?
	return typ, ok
}

func (s *SymbolTable) GetType() Type {
	return nil
}

func (s *SymbolTable) SymbolTypeName() string {
	return "symbol-table"
}

// Register will register the given symbol in this stab. If a
// symbol with the same name has alreayd been registered in this stab
// it will return false.
func (s *SymbolTable) Register(name string, sym SymbolValue) bool {
	if _, ok := s.Symbols[name]; ok {
		return false
	}
	s.Symbols[name] = sym
	return true
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

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Id:      rand.Intn(30000),
		Outer:   outer,
		Symbols: map[string]SymbolValue{},
		Types:   map[string]Type{},
	}
}
