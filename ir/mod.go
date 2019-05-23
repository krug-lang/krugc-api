package ir

import (
	"github.com/krug-lang/caasper/front"
)

// MODULE

// Module represents an individual krug file. contains functions,
// structures, impls, and the global statements (as a block).
type Module struct {
	Name           string                `json:"name,omitempty"`
	Structures     map[string]*Structure `json:"structures,omitempty"`
	StructureOrder []front.Token         `json:"structure_order,omitempty"`
	Functions      map[string]*Function  `json:"functions,omitempty"`
	FunctionOrder  []front.Token         `json:"function_order,omitempty"`
	Impls          map[string]*Impl      `json:"impls,omitempty"`
	ImplsOrder     []front.Token         `json:"impls_order,omitempty"`
	Global         *Block                `json:"global,omitempty"`
}

// NewModule creates a new module with the given name
func NewModule(name string) *Module {
	return &Module{
		name,

		map[string]*Structure{},
		[]front.Token{},

		map[string]*Function{},
		[]front.Token{},

		map[string]*Impl{},
		[]front.Token{},

		NewBlock(),
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
