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
	CharacterConstant                  = "charConst"
)

type ConstantNode struct {
	Kind                  ConstantNodeType
	VariableReferenceNode *VariableReferenceNode `json:"variableRefConst,omitempty"`
	IntegerConstantNode   *IntegerConstantNode   `json:"integerConst,omitempty"`
	FloatingConstantNode  *FloatingConstantNode  `json:"floatingConst,omitempty"`
	StringConstantNode    *StringConstantNode    `json:"stringConst,omitempty"`
	CharacterConstantNode *CharacterConstantNode `json:"charConst,omitempty"`
}

type VariableReferenceNode struct {
	Name Token `json:"name"`
}

type IntegerConstantNode struct {
	Value *big.Int `json:"value"`
}

type FloatingConstantNode struct {
	Value float64 `json:"value"`
}

type StringConstantNode struct {
	Value string `json:"value"`
}

type CharacterConstantNode struct {
	Value string `json:"value"`
}
