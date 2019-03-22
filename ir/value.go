package ir

import (
	"encoding/gob"
	"math/big"

	"github.com/krug-lang/server/front"
)

func init() {
	gob.Register(&IntegerValue{})
	gob.Register(&FloatingValue{})
	gob.Register(&StringValue{})
	gob.Register(&BinaryExpression{})
	gob.Register(&Identifier{})
	gob.Register(&Grouping{})
	gob.Register(&Builtin{})
	gob.Register(&UnaryExpression{})
	gob.Register(&Call{})
	gob.Register(&Path{})
	gob.Register(&Index{})
}

type Value interface {
	InferredType() Type
}

// PAREN EXPR

type Grouping struct {
	Val Value
}

func (g *Grouping) InferredType() Type {
	return g.Val.InferredType()
}

func NewGrouping(val Value) *Grouping {
	return &Grouping{val}
}

// BINARY

type BinaryExpression struct {
	LHand Value
	Op    string
	RHand Value
}

func (b *BinaryExpression) InferredType() Type {
	// TODO pick widest type, for now we pick the left type
	return b.LHand.InferredType()
}

func NewBinaryExpression(lh Value, op string, rh Value) *BinaryExpression {
	return &BinaryExpression{lh, op, rh}
}

// UNARY

type UnaryExpression struct {
	Op  string
	Val Value
}

func (u *UnaryExpression) InferredType() Type {
	// TODO?
	return u.Val.InferredType()
}

func NewUnaryExpression(op string, val Value) *UnaryExpression {
	return &UnaryExpression{op, val}
}

// FLOATING VALUE

type FloatingValue struct {
	Value float64
}

func (i *FloatingValue) InferredType() Type {
	return Float64
}

func NewFloatingValue(val float64) *FloatingValue {
	return &FloatingValue{val}
}

// INTEGER VALUE

type IntegerValue struct {
	RawValue *big.Int
}

func (i *IntegerValue) InferredType() Type {
	// TODO
	return Int32
}

func NewIntegerValue(val *big.Int) *IntegerValue {
	return &IntegerValue{val}
}

// STRING VALUE

type StringValue struct {
	Value string
}

func (s *StringValue) InferredType() Type {
	// TODO rune*
	return NewPointerType(Int32)
}

func NewStringValue(val string) *StringValue {
	return &StringValue{val}
}

// IDENTIFIER

type Identifier struct {
	Name front.Token
}

func (i *Identifier) InferredType() Type {
	panic("we need to deal with this later after name resolution")
}

func NewIdentifier(name front.Token) *Identifier {
	return &Identifier{name}
}

// BUILTIN

type Builtin struct {
	Name string
	Type Type
}

func (b *Builtin) InferredType() Type {
	// HM
	return b.Type
}

func NewBuiltin(name string, typ Type) *Builtin {
	return &Builtin{name, typ}
}

// CALL

type Call struct {
	Left   Value
	Params []Value
}

func (c *Call) InferredType() Type {
	panic("same thing as identifier, needs to be inferred from c.Left's ReturnType ")
}

func NewCall(left Value, params []Value) *Call {
	return &Call{left, params}
}

// PATH

type Path struct {
	Values []Value
}

func (p *Path) InferredType() Type {
	panic("also needs to be inferred in a later stage in sema")
}

func NewPath(values []Value) *Path {
	return &Path{values}
}

// INDEX

type Index struct {
	Left Value
	Sub  Value
}

func (i *Index) InferredType() Type {
	panic("needs to be inferred from i.Left's Base Type!")
}

func NewIndex(left, sub Value) *Index {
	return &Index{left, sub}
}
