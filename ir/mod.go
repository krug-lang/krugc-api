package ir

import (
	"encoding/gob"

	"github.com/hugobrains/krug-serv/front"
)

func init() {
	gob.Register(&Module{})
}

// MODULE

type Module struct {
	Name string

	Structures     map[string]*Structure
	StructureOrder []front.Token

	Functions     map[string]*Function
	FunctionOrder []front.Token

	Impls      map[string]*Impl
	ImplsOrder []front.Token
}

func NewModule(name string) *Module {
	return &Module{
		name,
		map[string]*Structure{},
		[]front.Token{},
		map[string]*Function{},
		[]front.Token{},
		map[string]*Impl{},
		[]front.Token{},
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
	if _, ok := m.Impls[i.Name.Value]; ok {
		return true
	}
	m.ImplsOrder = append(m.ImplsOrder, i.Name)
	m.Impls[i.Name.Value] = i
	return false
}

func (m *Module) RegisterStructure(s *Structure) {
	m.StructureOrder = append(m.StructureOrder, s.Name)
	m.Structures[s.Name.Value] = s
}

func (m *Module) RegisterFunction(f *Function) {
	m.FunctionOrder = append(m.FunctionOrder, f.Name)
	m.Functions[f.Name.Value] = f
}
