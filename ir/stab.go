package ir

import (
	"encoding/gob"
	"math/rand"
)

func init() {
	gob.Register(&SymbolTable{})
	gob.Register(&Symbol{})
}

type SymbolValue interface{}

type Symbol struct {
	Name string
}

func NewSymbol(name string) *Symbol {
	return &Symbol{name}
}

type SymbolTable struct {
	Id      int
	Outer   *SymbolTable
	Symbols map[string]SymbolValue
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
	}
}
