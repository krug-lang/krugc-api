package front

// TypeNodeType ...
type TypeNodeType string

// ...
const (
	UnresolvedType TypeNodeType = "unresolvedType"
	PointerType                 = "pointerType"
	ArrayType                   = "arrayType"
	TupleType                   = "tupleType"
	StructureType               = "structType"
)

// UnresolvedTypeNode ...
type UnresolvedTypeNode struct {
	Name     string
	Resolved bool
}

// PointerTypeNode ...
// *Type
type PointerTypeNode struct {
	Base *ExpressionNode
}

// TupleTypeNode ...
type TupleTypeNode struct {
	Types []*ExpressionNode
}

// ArrayTypeNode ...
// [TypeNode; expr]
type ArrayTypeNode struct {
	Base *ExpressionNode
	Size *ExpressionNode
}

// StructureTypeNode ...
// "struct" iden { ... }
type StructureTypeNode struct {
	Name   Token        `json:"name"`
	Fields []*NamedType `json:"fields"`
}

// TypeNode ...
type TypeNode struct {
	Kind TypeNodeType

	TupleTypeNode      *TupleTypeNode      `json:"tupleType,omitempty"`
	UnresolvedTypeNode *UnresolvedTypeNode `json:"unresolvedType,omitempty"`
	PointerTypeNode    *PointerTypeNode    `json:"pointerType,omitempty"`
	ArrayTypeNode      *ArrayTypeNode      `json:"arrayType,omitempty"`
	StructureTypeNode  *StructureTypeNode  `json:"structType,omitempty"`
}
