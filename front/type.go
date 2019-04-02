package front

// lol
type TypeNodeType int

const (
	UnresolvedType TypeNodeType = iota
	PointerType
	ArrayType
)

// SomeStructureName
type UnresolvedTypeNode struct {
	Name string
}

// *TypeNode
type PointerTypeNode struct {
	Base *TypeNode
}

// [TypeNode]
type ArrayTypeNode struct {
	Base *TypeNode
	Size *ExpressionNode
}

type TypeNode struct {
	Kind TypeNodeType

	UnresolvedTypeNode *UnresolvedTypeNode
	PointerTypeNode    *PointerTypeNode
	ArrayTypeNode      *ArrayTypeNode
}
