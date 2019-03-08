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
}

type ExpressionNode interface {
	Print() string
}

type VariableNode struct {
	Name string
}

func (v *VariableNode) Print() string {
	return fmt.Sprintf("(v %s)", v.Name)
}

func NewVariable(name string) *VariableNode {
	return &VariableNode{name}
}

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

type UnaryExpression struct {
	Operator string
	Value    ExpressionNode
}

func (u *UnaryExpression) Print() string {
	return fmt.Sprintf("(u %s%s)", u.Operator, u.Value.Print())
}

func NewUnary(op string, value ExpressionNode) *UnaryExpression {
	return &UnaryExpression{op, value}
}

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

type Grouping struct {
	Value ExpressionNode
}

func (g *Grouping) Print() string {
	return fmt.Sprintf("(g %s)", g.Value.Print())
}

func NewGrouping(val ExpressionNode) *Grouping {
	return &Grouping{val}
}
