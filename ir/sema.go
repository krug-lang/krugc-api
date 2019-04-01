package ir

type SemanticModule struct {
	ScopeMap *ScopeMap
	TypeMap  *TypeMap
	Module   *Module
}

func NewSemanticModule(sm *ScopeMap, tm *TypeMap, m *Module) SemanticModule {
	return SemanticModule{sm, tm, m}
}

// ScopeMap is a data structure that
// contains the scope information for a single
// module.
//
// the scope map is purely for the symbols in
// data, it does not contain any type information.
type ScopeMap struct {
	// TODO
	Functions  map[string]*SymbolTable
	Structures map[string]*SymbolTable
}

func (s *ScopeMap) RegisterFunction(name string, sym *SymbolTable) bool {
	if _, ok := s.Functions[name]; ok {
		return ok
	}
	s.Functions[name] = sym
	return true
}

func (s *ScopeMap) RegisterStructure(name string, sym *SymbolTable) bool {
	if _, ok := s.Structures[name]; ok {
		return ok
	}
	s.Structures[name] = sym
	return true
}

func NewScopeMap() *ScopeMap {
	return &ScopeMap{
		Functions:  map[string]*SymbolTable{},
		Structures: map[string]*SymbolTable{},
	}
}

// the type map data structure is a per module
// data structure that contains information of
// all of the symbols and their types.
type TypeMap struct {
}

func NewTypeMap() *TypeMap {
	return &TypeMap{}
}
