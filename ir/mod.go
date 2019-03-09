package ir

import "encoding/gob"

func init() {
	gob.Register(&Structure{})
	gob.Register(&Function{})
	gob.Register(&Module{})
}

type Structure struct {
	Name    string
	Fields  map[string]Type
	Methods map[string]Function
}

func NewStructure(name string, fields map[string]Type) *Structure {
	return &Structure{name, fields, map[string]Function{}}
}

func (s *Structure) RegisterMethod(f Function) {
	s.Methods[f.Name] = f
}

type Function struct {
	Name       string
	Param      map[string]Type
	ReturnType Type
	Body       *Block
}

func NewFunction(name string, params map[string]Type, ret Type) *Function {
	return &Function{name, params, ret, NewBlock()}
}

type Module struct {
	Name string

	Structures map[string]*Structure
	Functions  map[string]*Function
}

func NewModule(name string) *Module {
	return &Module{
		name,
		map[string]*Structure{},
		map[string]*Function{},
	}
}

func (m *Module) GetStructure(name string) (*Structure, bool) {
	s, ok := m.Structures[name]
	return s, ok
}

func (m *Module) RegisterStructure(s *Structure) {
	m.Structures[s.Name] = s
}

func (m *Module) RegisterFunction(f *Function) {
	m.Functions[f.Name] = f
}
