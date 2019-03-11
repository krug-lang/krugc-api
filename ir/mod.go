package ir

import "encoding/gob"

func init() {
	gob.Register(&Module{})
}

// MODULE

type Module struct {
	Name       string
	Root       *SymbolTable
	Structures map[string]*Structure
	Functions  map[string]*Function
	Impls      map[string]*Impl
}

func NewModule(name string) *Module {
	return &Module{
		name,
		nil,
		map[string]*Structure{},
		map[string]*Function{},
		map[string]*Impl{},
	}
}

func (m *Module) GetStructure(name string) (*Structure, bool) {
	s, ok := m.Structures[name]
	return s, ok
}

func (m *Module) GetImpl(name string) (*Impl, bool) {
	s, ok := m.Impls[name]
	return s, ok
}

// RegisterImpl registers the given impl with the module.
// returns TRUE if an impl has already been registered with
// the same name and FALSE if not.
func (m *Module) RegisterImpl(i *Impl) (failed bool) {
	if _, ok := m.Impls[i.Name]; ok {
		return true
	}
	m.Impls[i.Name] = i
	return false
}

func (m *Module) RegisterStructure(s *Structure) {
	m.Structures[s.Name] = s
}

func (m *Module) RegisterFunction(f *Function) {
	m.Functions[f.Name] = f
}
