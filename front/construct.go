package front

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&IfStatement{})
	gob.Register(&ElseIfStatement{})
	gob.Register(&LoopStatement{})
	gob.Register(&WhileLoopStatement{})
}

type WhileLoopStatement struct {
	Cond  ExpressionNode
	Post  ExpressionNode
	Block []StatementNode
}

func (w *WhileLoopStatement) Print() string {
	return fmt.Sprintf("while %s ... ", w.Cond.Print())
}

func NewWhileLoopStatement(cond ExpressionNode, post ExpressionNode, block []StatementNode) *WhileLoopStatement {
	return &WhileLoopStatement{cond, post, block}
}

type LoopStatement struct {
	Block []StatementNode
}

func (l *LoopStatement) Print() string {
	// TODO
	return fmt.Sprintf("loop ...")
}

func NewLoopStatement(block []StatementNode) *LoopStatement {
	return &LoopStatement{block}
}

type ElseIfStatement struct {
	Cond  ExpressionNode
	Block []StatementNode
}

func NewElseIfStatement(cond ExpressionNode, block []StatementNode) *ElseIfStatement {
	return &ElseIfStatement{cond, block}
}

type IfStatement struct {
	Cond    ExpressionNode
	Block   []StatementNode
	ElseIfs []*ElseIfStatement
	Else    []StatementNode
}

func NewIfStatement(cond ExpressionNode, block []StatementNode, elseIfs []*ElseIfStatement, elseBlock []StatementNode) *IfStatement {
	return &IfStatement{cond, block, elseIfs, elseBlock}
}

func (i *IfStatement) Print() string {
	return fmt.Sprintf("if %s", i.Cond)
}
