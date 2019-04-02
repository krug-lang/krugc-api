package front

import (
	"math/big"
)

type ConstantNodeType int

const (
	VariableReference ConstantNodeType = iota
	IntegerConstant
	FloatingConstant
	StringConstant
)

type ConstantNode struct {
	Kind                  ConstantNodeType
	VariableReferenceNode *VariableReferenceNode
	IntegerConstantNode   *IntegerConstantNode
	FloatingConstantNode  *FloatingConstantNode
	StringConstantNode    *StringConstantNode
}

type VariableReferenceNode struct {
	Name Token
}

type IntegerConstantNode struct {
	Value *big.Int
}

type FloatingConstantNode struct {
	Value float64
}

type StringConstantNode struct {
	Value string
}
