package ir

import (
	"github.com/hugobrains/caasper/front"
)

type InstructionKind string

const (
	BlockInstr           = "blockInstr"
	AssignInstr          = "assignInstr"
	LocalInstr           = "localInstr"
	AllocaInstr          = "allocaInstr"
	NextInstr            = "nextInstr"
	BreakInstr           = "breakInstr"
	ReturnInstr          = "returnInstr"
	LoopInstr            = "loopInstr"
	WhileLoopInstr       = "whileLoopInstr"
	ElseIfStatementInstr = "elseIfStatementInstr"
	IfStatementInstr     = "ifStatementInstr"
	ExpressionInstr      = "exprInstr"
	DeferInstr           = "deferInstr"
)

type Instruction struct {
	Kind                InstructionKind
	Block               *Block
	Assign              *Assign
	Local               *Local
	Alloca              *Alloca
	Next                *Next
	Break               *Break
	Return              *Return
	Loop                *Loop
	Defer               *Defer
	WhileLoop           *WhileLoop
	ElseIfStatement     *ElseIfStatement
	IfStatement         *IfStatement
	ExpressionStatement *Value
}

// BLOCK

type Block struct {
	DeferStack []*Defer
	Instr      []*Instruction
	Stab       *SymbolTable
	Return     *Instruction
}

func (b *Block) PushDefer(def *Defer) {
	b.DeferStack = append(b.DeferStack, def)
}

func (b *Block) SetReturn(ret *Instruction) {
	b.Return = ret
}

func (b *Block) AddInstr(instr *Instruction) {
	b.Instr = append(b.Instr, instr)
}

func NewBlock() *Block {
	return &Block{
		[]*Defer{},
		[]*Instruction{},
		nil,
		nil,
	}
}

// ASSIGN

type Assign struct {
	LHand *Value
	Op    string
	RHand *Value
}

func (a *Assign) InferredType() Type {
	panic("handle this in sema?")
}

func NewAssign(lh *Value, op string, rh *Value) *Assign {
	return &Assign{lh, op, rh}
}

type Defer struct {
	Stat  *Instruction
	Block *Block
}

func NewDefer(stat *Instruction, block *Block) *Defer {
	return &Defer{stat, block}
}

// LOCAL

type Local struct {
	Name    front.Token
	Type    *Type
	Mutable bool
	Val     *Value
}

func (l *Local) SetValue(v *Value) {
	l.Val = v
}

func (l *Local) SetMutable(m bool) {
	l.Mutable = m
}

func NewLocal(name front.Token, typ *Type) *Local {
	return &Local{name, typ, false, nil}
}

// ALLOCA

type Alloca struct {
	Name front.Token
	Type *Type
	Val  *Value
}

func (a *Alloca) SetValue(v *Value) {
	a.Val = v
}

func NewAlloca(name front.Token, typ *Type) *Alloca {
	return &Alloca{name, typ, nil}
}

// NEXT, BREAK

type Next struct{}

func NewNext() *Next { return &Next{} }

type Break struct{}

func NewBreak() *Break { return &Break{} }

// RETURN

type Return struct {
	Val *Value
}

func NewReturn(val *Value) *Return {
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
	Cond *Value
	Post *Value
	Body *Block
}

func NewWhileLoop(cond *Value, post *Value, body *Block) *WhileLoop {
	return &WhileLoop{cond, post, body}
}

// ELSE IF

type ElseIfStatement struct {
	Cond *Value
	Body *Block
}

func NewElseIfStatement(cond *Value, body *Block) *ElseIfStatement {
	return &ElseIfStatement{cond, body}
}

// IF STATEMENT

type IfStatement struct {
	Cond   *Value
	True   *Block
	ElseIf []*ElseIfStatement
	Else   *Block
}

func NewIfStatement(cond *Value, t *Block, elses []*ElseIfStatement, f *Block) *IfStatement {
	return &IfStatement{cond, t, elses, f}
}
