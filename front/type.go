package front

type TypeNodeType int

const (
	UnresolvedType TypeNodeType = iota
	PointerType
	ArrayType
	TupleType
)

// SomeStructureName
type UnresolvedTypeNode struct {
	Name string
}

// *TypeNode
type PointerTypeNode struct {
	Base *TypeNode
}

type TupleTypeNode struct {
	Types []*TypeNode
}

// [TypeNode]
type ArrayTypeNode struct {
	Base *TypeNode
	Size *ExpressionNode
}

type TypeNode struct {
	Kind TypeNodeType

	TupleTypeNode      *TupleTypeNode      `json:"tupleType,omitempty"`
	UnresolvedTypeNode *UnresolvedTypeNode `json:"unresolvedType,omitempty"`
	PointerTypeNode    *PointerTypeNode    `json:"pointerType,omitempty"`
	ArrayTypeNode      *ArrayTypeNode      `json:"arrayType,omitempty"`
}
