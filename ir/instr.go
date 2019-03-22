package ir

import (
	"encoding/gob"

	"github.com/krug-lang/server/front"
)

func init() {
	gob.Register(&Block{})
	gob.Register(&Local{})
	gob.Register(&Alloca{})
	gob.Register(&Return{})
	gob.Register(&Break{})
	gob.Register(&Next{})
	gob.Register(&Loop{})
	gob.Register(&WhileLoop{})
	gob.Register(&IfStatement{})
	gob.Register(&Assign{})
	gob.Register(&ElseIfStatement{})
}

type Instruction interface{}

// BLOCK

type Block struct {
	Instr []Instruction
	Stab  *SymbolTable
}

func (b *Block) AddInstr(instr Instruction) {
	b.Instr = append(b.Instr, instr)
}

func NewBlock() *Block {
	return &Block{
		[]Instruction{},
		nil,
	}
}

// ASSIGN

type Assign struct {
	LHand Value
	Op    string
	RHand Value
}

func (a *Assign) InferredType() Type {
	panic("handle this in sema?")
}

func NewAssign(lh Value, op string, rh Value) *Assign {
	return &Assign{lh, op, rh}
}

// LOCAL

type Local struct {
	Name    front.Token
	Type    Type
	Mutable bool
	Val     Value
}

func (l *Local) SetValue(v Value) {
	l.Val = v
}

func (l *Local) SetMutable(m bool) {
	l.Mutable = m
}

func NewLocal(name front.Token, typ Type) *Local {
	return &Local{name, typ, false, nil}
}

// ALLOCA

type Alloca struct {
	Name front.Token
	Type Type
	Val  Value
}

func (a *Alloca) SetValue(v Value) {
	a.Val = v
}

func NewAlloca(name front.Token, typ Type) *Alloca {
	return &Alloca{name, typ, nil}
}

// NEXT, BREAK

type Next struct{}

func NewNext() *Next { return &Next{} }

type Break struct{}

func NewBreak() *Break { return &Break{} }

// RETURN

type Return struct {
	Val Value
}

func NewReturn(val Value) *Return {
	return &Return{val}
}

// LOOP

type Loop struct {
	Body *Block
}

func NewLoop(body *Block) *Loop {
	return &Loop{body}
}

// WHILE LOOP

type WhileLoop struct {
	Cond Value
	Post Value
	Body *Block
}

func NewWhileLoop(cond Value, post Value, body *Block) *WhileLoop {
	return &WhileLoop{cond, post, body}
}

// ELSE IF

type ElseIfStatement struct {
	Cond Value
	Body *Block
}

func NewElseIfStatement(cond Value, body *Block) *ElseIfStatement {
	return &ElseIfStatement{cond, body}
}

// IF STATEMENT

type IfStatement struct {
	Cond   Value
	True   *Block
	ElseIf []*ElseIfStatement
	Else   *Block
}

func NewIfStatement(cond Value, t *Block, elses []*ElseIfStatement, f *Block) *IfStatement {
	return &IfStatement{cond, t, elses, f}
}
