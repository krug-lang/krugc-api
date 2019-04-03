package ir

import (
	"fmt"
	"math/big"

	"github.com/hugobrains/caasper/front"
)

type ValueKind string

const (
	// fix the ValueValue this is from a copy paste job.
	IntegerValueValue     = "IntegerValue"
	FloatingValueValue    = "FloatingValue"
	StringValueValue      = "StringValue"
	BinaryExpressionValue = "BinaryExpression"
	IdentifierValue       = "Identifier"
	BuiltinValue          = "Builtin"
	GroupingValue         = "Grouping"
	UnaryExpressionValue  = "UnaryExpression"
	CallValue             = "Call"
	PathValue             = "Path"
	IndexValue            = "Index"
	AssignValue           = "Assign"
)

type Value struct {
	Kind             ValueKind
	IntegerValue     *IntegerValue
	FloatingValue    *FloatingValue
	StringValue      *StringValue
	BinaryExpression *BinaryExpression
	Identifier       *Identifier
	Grouping         *Grouping
	Assign           *Assign
	Builtin          *Builtin
	UnaryExpression  *UnaryExpression
	Call             *Call
	Path             *Path
	Index            *Index
}

// FIXME! this is shit
func (v *Value) InferredType() *Type {
	switch v.Kind {
	case IntegerValueValue:
		return v.IntegerValue.InferredType()
	case FloatingValueValue:
		return v.FloatingValue.InferredType()
	case StringValueValue:
		return v.StringValue.InferredType()
	case BinaryExpressionValue:
		return v.BinaryExpression.InferredType()
	case IdentifierValue:
		return v.Identifier.InferredType()
	case BuiltinValue:
		return v.Builtin.InferredType()
	case GroupingValue:
		return v.Grouping.InferredType()
	case UnaryExpressionValue:
		return v.UnaryExpression.InferredType()
	case CallValue:
		return v.Call.InferredType()
	case PathValue:
		return v.Path.InferredType()
	case IndexValue:
		return v.Index.InferredType()
	case AssignValue:
		panic("uh")
	default:
		panic("unhandled Value::InferredType()")
	}
}

// PAREN EXPR

type Grouping struct {
	Val *Value
}

func (g *Grouping) InferredType() *Type {
	return g.Val.InferredType()
}

func NewGrouping(val *Value) *Grouping {
	return &Grouping{val}
}

// BINARY

type BinaryExpression struct {
	LHand *Value
	Op    string
	RHand *Value
}

func (b *BinaryExpression) InferredType() *Type {
	// TODO pick widest type, for now we pick the left type
	fmt.Println(b.LHand, b.Op, b.RHand)
	return b.LHand.InferredType()
}

func NewBinaryExpression(lh *Value, op string, rh *Value) *BinaryExpression {
	return &BinaryExpression{lh, op, rh}
}

// UNARY

type UnaryExpression struct {
	Op  string
	Val *Value
}

func (u *UnaryExpression) InferredType() *Type {
	// TODO?
	return u.Val.InferredType()
}

func NewUnaryExpression(op string, val *Value) *UnaryExpression {
	return &UnaryExpression{op, val}
}

// FLOATING VALUE

type FloatingValue struct {
	Value float64
}

func (i *FloatingValue) InferredType() *Type {
	return Float64
}

func NewFloatingValue(val float64) *FloatingValue {
	return &FloatingValue{val}
}

// INTEGER VALUE

type IntegerValue struct {
	RawValue *big.Int
}

func (i *IntegerValue) InferredType() *Type {
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

func (s *StringValue) InferredType() *Type {
	// TODO rune*
	return &Type{
		Kind:    PointerKind,
		Pointer: NewPointerType(Int32),
	}
}

func NewStringValue(val string) *StringValue {
	return &StringValue{val}
}

// IDENTIFIER

type Identifier struct {
	Name front.Token
}

func (i *Identifier) InferredType() *Type {
	panic("we need to deal with this later after name resolution")
}

func NewIdentifier(name front.Token) *Identifier {
	return &Identifier{name}
}

// BUILTIN

type Builtin struct {
	Name string
	Type *Type
}

func (b *Builtin) InferredType() *Type {
	// HM
	return b.Type
}

func NewBuiltin(name string, typ *Type) *Builtin {
	return &Builtin{name, typ}
}

// CALL

type Call struct {
	Left   *Value
	Params []*Value
}

func (c *Call) InferredType() *Type {
	panic("same thing as identifier, needs to be inferred from c.Left's ReturnType ")
}

func NewCall(left *Value, params []*Value) *Call {
	return &Call{left, params}
}

// PATH

type Path struct {
	Values []*Value
}

func (p *Path) InferredType() *Type {
	panic("also needs to be inferred in a later stage in sema")
}

func NewPath(values []*Value) *Path {
	return &Path{values}
}

// INDEX

type Index struct {
	Left *Value
	Sub  *Value
}

func (i *Index) InferredType() *Type {
	panic("needs to be inferred from i.Left's Base Type!")
}

func NewIndex(left, sub *Value) *Index {
	return &Index{left, sub}
}
