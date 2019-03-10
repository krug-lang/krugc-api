package front

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(ParseTree{})
	gob.Register(&StructureDeclaration{})
	gob.Register(&ImplDeclaration{})
	gob.Register(&TraitDeclaration{})
	gob.Register(&FunctionDeclaration{})
	gob.Register(&FunctionPrototypeDeclaration{})
}

type NamedType struct {
	Name string
	Type TypeNode
}

//
// STRUCT
//

// "struct" iden { ... }
type StructureDeclaration struct {
	Name   string
	Fields []*NamedType
}

func (s *StructureDeclaration) Print() string {
	return fmt.Sprintf("(struct %s, fields %s)", s.Name, s.Fields)
}

func NewStructureDeclaration(name string, fields []*NamedType) *StructureDeclaration {
	return &StructureDeclaration{
		Name:   name,
		Fields: fields,
	}
}

func (s *StructureDeclaration) NodeName() string {
	return "struct-decl"
}

//
// TRAIT
//

// "trait" iden { ... }
type TraitDeclaration struct {
	Name    string
	Members []*FunctionPrototypeDeclaration
}

func NewTraitDeclaration(name string, members []*FunctionPrototypeDeclaration) *TraitDeclaration {
	return &TraitDeclaration{
		name, members,
	}
}

func (t *TraitDeclaration) Print() string {
	return fmt.Sprintf("(trait %s)", t.Name)
}

func (t *TraitDeclaration) NodeName() string {
	return "trait-decl"
}

//
// IMPL
//

// impl Name
type ImplDeclaration struct {
	Name      string
	Functions []*FunctionDeclaration
}

func (i *ImplDeclaration) Print() string {
	return fmt.Sprintf("(impl %s)", i.Name)
}

func NewImplDeclaration(name string, funcs []*FunctionDeclaration) *ImplDeclaration {
	return &ImplDeclaration{
		name, funcs,
	}
}

func (i *ImplDeclaration) NodeName() string {
	return "impl-decl"
}

//
// FUNC PROTO
//

// "func" iden "(" args ")"
type FunctionPrototypeDeclaration struct {
	Name      string
	Arguments []*NamedType

	// TODO should this be set to anything by
	// default, e.g. we can inject a "void"
	// into here?
	ReturnType TypeNode
}

func (f *FunctionPrototypeDeclaration) SetReturnType(r TypeNode) {
	f.ReturnType = r
}

func (f *FunctionPrototypeDeclaration) Print() string {
	return fmt.Sprintf("(func %s, args %s)", f.Name, f.Arguments)
}

func NewFunctionPrototypeDeclaration(name string, args []*NamedType) *FunctionPrototypeDeclaration {
	return &FunctionPrototypeDeclaration{
		Name:      name,
		Arguments: args,
	}
}

func (f *FunctionPrototypeDeclaration) NodeName() string {
	return "func-proto-decl"
}

//
// FUNC DECL
//

// [ FunctionPrototypeDeclaration ] "{" { Stat ";" } "}"
type FunctionDeclaration struct {
	*FunctionPrototypeDeclaration
	Body []StatementNode
}

func (f *FunctionDeclaration) Print() string {
	var b string
	for _, st := range f.Body {
		b += fmt.Sprintf("\t%s\n", st.Print())
	}
	return fmt.Sprintf("(! %s) => \n(%s)", f.FunctionPrototypeDeclaration.Print(), b)
}

func NewFunctionDeclaration(proto *FunctionPrototypeDeclaration, body []StatementNode) *FunctionDeclaration {
	return &FunctionDeclaration{
		proto, body,
	}
}

func (f *FunctionDeclaration) NodeName() string {
	return "func-decl"
}
