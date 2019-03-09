package front

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&LiteralExpression{})
	gob.Register(&UnaryExpression{})
	gob.Register(&BinaryExpression{})
	gob.Register(&Grouping{})
	gob.Register(&VariableNode{})
	gob.Register(&BuiltinExpression{})
}

type ExpressionNode interface {
	Print() string
}

// BUILTIN

type BuiltinExpression struct {
	Name string
	Type TypeNode
}

func (b *BuiltinExpression) Print() string {
	return fmt.Sprintf("%s!%s", b.Name, b.Type)
}

func NewBuiltinExpression(name string, typ TypeNode) *BuiltinExpression {
	return &BuiltinExpression{name, typ}
}

// VARIABLE

type VariableNode struct {
	Name string
}

func (v *VariableNode) Print() string {
	return fmt.Sprintf("(v %s)", v.Name)
}

func NewVariable(name string) *VariableNode {
	return &VariableNode{name}
}

// LITERAL

// string, number
type LiteralExpression struct {
	Value string
}

func (l *LiteralExpression) Print() string {
	return fmt.Sprintf("(lit %s)", l.Value)
}

func NewLiteral(value string) *LiteralExpression {
	return &LiteralExpression{value}
}

// UNARY

type UnaryExpression struct {
	Operator string
	Value    ExpressionNode
}

func (u *UnaryExpression) Print() string {
	return fmt.Sprintf("(u %s%s)", u.Operator, u.Value.Print())
}

func NewUnaryExpression(op string, value ExpressionNode) *UnaryExpression {
	return &UnaryExpression{op, value}
}

// BINARY

type BinaryExpression struct {
	LHand    ExpressionNode
	Operator string
	RHand    ExpressionNode
}

func (b *BinaryExpression) Print() string {
	return fmt.Sprintf("(%s%s%s)", b.LHand.Print(), b.Operator, b.RHand.Print())
}

func NewBinary(lh ExpressionNode, op string, rh ExpressionNode) *BinaryExpression {
	return &BinaryExpression{lh, op, rh}
}

// GROUPING

type Grouping struct {
	Value ExpressionNode
}

func (g *Grouping) Print() string {
	return fmt.Sprintf("(g %s)", g.Value.Print())
}

func NewGrouping(val ExpressionNode) *Grouping {
	return &Grouping{val}
}
