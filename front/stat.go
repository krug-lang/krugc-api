package front

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&LetStatement{})
	gob.Register(&MutableStatement{})
	gob.Register(&ReturnStatement{})
}

type StatementNode interface {
	Print() string
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
	Name  string
	Type  TypeNode
	Value ExpressionNode
}

func NewLetStatement(name string, kind TypeNode, val ExpressionNode) *LetStatement {
	return &LetStatement{name, kind, val}
}

func (l *LetStatement) Print() string {
	return fmt.Sprintf("let %s = ", l.Name)
}

// "mut" iden [ ":" Type ] [ = Value ];
type MutableStatement struct {
	Name  string
	Type  TypeNode
	Value ExpressionNode
}

func NewMutableStatement(name string, typ TypeNode, val ExpressionNode) *MutableStatement {
	return &MutableStatement{name, typ, val}
}

func (m *MutableStatement) Print() string {
	return fmt.Sprintf("mut %s = ", m.Name)
}
