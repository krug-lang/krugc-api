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
