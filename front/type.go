package front

import (
	"encoding/gob"
	"fmt"
)

/*
	Array types should have a length associated
	with them.
*/

type TypeNode interface{}

func init() {
	gob.Register(&UnresolvedType{})
	gob.Register(&PointerType{})
	gob.Register(&ArrayType{})
}

// SomeStructureName
type UnresolvedType struct {
	Name string
}

func (u *UnresolvedType) String() string {
	return fmt.Sprintf("(? %s)", u.Name)
}

func NewUnresolvedType(name string) *UnresolvedType {
	return &UnresolvedType{
		name,
	}
}

// *TypeNode
type PointerType struct {
	Base TypeNode
}

func (p *PointerType) String() string {
	return fmt.Sprintf("(ptr %s)", p.Base)
}

func NewPointerType(base TypeNode) *PointerType {
	return &PointerType{
		base,
	}
}

// [TypeNode]
type ArrayType struct {
	Base TypeNode
	Size ExpressionNode
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("(arr %s)", a.Base)
}

func NewArrayType(base TypeNode, size ExpressionNode) *ArrayType {
	return &ArrayType{
		base, size,
	}
}
