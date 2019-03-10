package front

import (
	"encoding/gob"
	"fmt"
	"math/big"
	"strconv"
)

func init() {
	gob.Register(&IntegerConstant{})
	gob.Register(&FloatingConstant{})
	gob.Register(&VariableReference{})
	gob.Register(&StringConstant{})
}

type VariableReference struct {
	Name string
}

func (v *VariableReference) Print() string {
	return v.Name
}

func NewVariableReference(name string) *VariableReference {
	return &VariableReference{name}
}

type ConstantNode interface{}

// INTEGER CONST

type IntegerConstant struct {
	Value *big.Int
}

func (i *IntegerConstant) Print() string {
	return fmt.Sprintf("(ic %s)", i.Value)
}

func NewIntegerConstant(val string) *IntegerConstant {
	value := new(big.Int)
	value, ok := value.SetString(val, 10)
	if !ok {
		panic("bad int const")
	}
	return &IntegerConstant{value}
}

// FLOAT CONST

type FloatingConstant struct {
	Value float64
}

func (f *FloatingConstant) Print() string {
	return fmt.Sprintf("(fc %f)", f.Value)
}

func NewFloatingConstant(val string) *FloatingConstant {
	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return &FloatingConstant{value}
}

// STRING CONST

type StringConstant struct {
	Value string
}

func (s *StringConstant) Print() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

func NewStringConstant(val string) *StringConstant {
	return &StringConstant{val}
}
