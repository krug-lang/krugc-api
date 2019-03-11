package front

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&LetStatement{})
	gob.Register(&MutableStatement{})
	gob.Register(&ReturnStatement{})
	gob.Register(&AssignStatement{})
	gob.Register(&BreakStatement{})
	gob.Register(&NextStatement{})
	gob.Register(&BlockNode{})
}

type StatementNode interface {
	Print() string
}

type NextStatement struct{}

func (n *NextStatement) Print() string {
	return "next"
}

func NewNextStatement() *NextStatement {
	return &NextStatement{}
}

type BreakStatement struct{}

func (b *BreakStatement) Print() string {
	return "break"
}

func NewBreakStatement() *BreakStatement {
	return &BreakStatement{}
}

type ReturnStatement struct {
	Value ExpressionNode
}

func (r *ReturnStatement) Print() string {
	return fmt.Sprintf("ret %s", r.Value)
}

func NewReturnStatement(val ExpressionNode) *ReturnStatement {
	return &ReturnStatement{val}
}

// "let" iden [ ":" Type ] = Value;
type LetStatement struct {
	Name  Token
	Type  TypeNode
	Value ExpressionNode
}

func NewLetStatement(name Token, kind TypeNode, val ExpressionNode) *LetStatement {
	return &LetStatement{name, kind, val}
}

func (l *LetStatement) Print() string {
	return fmt.Sprintf("let %s = ", l.Name)
}

// "mut" iden [ ":" Type ] [ = Value ];
type MutableStatement struct {
	Name  Token
	Type  TypeNode
	Value ExpressionNode
}

func NewMutableStatement(name Token, typ TypeNode, val ExpressionNode) *MutableStatement {
	return &MutableStatement{name, typ, val}
}

func (m *MutableStatement) Print() string {
	return fmt.Sprintf("mut %s = ", m.Name)
}

// BLOCK NODE

type BlockNode struct {
	Stats []StatementNode
}

func (b *BlockNode) Print() string {
	return "{ ... }"
}

func NewBlockNode(stats []StatementNode) *BlockNode {
	return &BlockNode{stats}
}

// ASSIGN

type AssignStatement struct {
	LHand ExpressionNode
	Op    string
	RHand ExpressionNode
}

func (a *AssignStatement) Print() string {
	return fmt.Sprintf("%s %s %s", a.LHand.Print(), a.Op, a.RHand.Print())
}

func NewAssignmentStatement(lh ExpressionNode, op string, rh ExpressionNode) *AssignStatement {
	return &AssignStatement{lh, op, rh}
}
