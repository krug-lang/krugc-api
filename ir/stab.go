package ir

import (
	"fmt"
	"math/rand"

	"github.com/krug-lang/caasper/front"
)

type SymbolValueKind string

const (
	SymbolKind      SymbolValueKind = "sym"
	SymbolTableKind                 = "stab"
)

type SymbolValue struct {
	Kind        SymbolValueKind
	Symbol      *Symbol
	SymbolTable *SymbolTable
}

type Symbol struct {
	Name    front.Token `json:"name"`
	Owned   bool        `json:"owned"`
	Mutable bool        `json:"mutable"`
}

func (s *Symbol) GetType() *Type {
	return nil
}

func (s *Symbol) SymbolTypeName() string {
	return "symbol"
}

func NewSymbol(name front.Token, owned bool, mutable bool) *Symbol {
	return &Symbol{name, owned, mutable}
}

type SymbolTable struct {
	ID        int                     `json:"id"`
	OuterID   int                     `json:"outer_id,omitempty"`
	Inner     []*SymbolTable          `json:"inner,omitempty"`
	Types     map[string]*Type        `json:"types"`
	Symbols   map[string]*SymbolValue `json:"symbols"`
	SymbolSet []string                `json:"symbol_set,omitempty"`
}

func (s *SymbolTable) String() string {
	res := "{"

	idx := 0
	for _, sym := range s.Symbols {
		if idx != 0 {
			res += " "
		}

		// FIXMe
		switch sym.Kind {
		case SymbolKind:
			res += sym.Symbol.Name.Value
		}

		idx++
	}
	res += "}"
	return res
}

func (s *SymbolTable) RegisterType(name string, t *Type) {
	s.Types[name] = t
}

func (s *SymbolTable) LookupType(name string) (*Type, bool) {
	typ, ok := s.Types[name]
	return typ, ok
}

func (s *SymbolTable) GetType() *Type {
	return nil
}

func (s *SymbolTable) SymbolTypeName() string {
	return "symbol-table"
}

// Register will register the given symbol in this stab. If a
// symbol with the same name has alreayd been registered in this stab
// it will return false.
func (s *SymbolTable) Register(name string, sym *SymbolValue) bool {
	if _, ok := s.Symbols[name]; ok {
		return false
	}
	s.SymbolSet = append(s.SymbolSet, name)
	s.Symbols[name] = sym
	return true
}

func (s *SymbolTable) Lookup(name string) (*SymbolValue, bool) {
	if sym, ok := s.Symbols[name]; ok {
		if ok {
			return sym, ok
		}
	}

	panic(fmt.Sprintf("coun't find symbol %s", name))
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	outerID := -1
	if outer != nil {
		outerID = outer.ID
	}
	return &SymbolTable{
		ID:        rand.Intn(30000),
		OuterID:   outerID,
		Inner:     []*SymbolTable{},
		Symbols:   map[string]*SymbolValue{},
		Types:     map[string]*Type{},
		SymbolSet: []string{},
	}
}
