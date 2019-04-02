package front

import (
	"math/big"
)

type ConstantNodeType string

const (
	VariableReference ConstantNodeType = "variableRefConst"
	IntegerConstant                    = "integerConst"
	FloatingConstant                   = "floatingConst"
	StringConstant                     = "stringConst"
)

type ConstantNode struct {
	Kind                  ConstantNodeType
	VariableReferenceNode *VariableReferenceNode `json:"variableRefConst,omitempty"`
	IntegerConstantNode   *IntegerConstantNode   `json:"integerConst,omitempty"`
	FloatingConstantNode  *FloatingConstantNode  `json:"floatingConst,omitempty"`
	StringConstantNode    *StringConstantNode    `json:"stringConst,omitempty"`
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
