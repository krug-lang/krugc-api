package ir

import (
	"strings"

	"github.com/krug-lang/caasper/front"
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
	JumpInstr            = "jumpInstr"
	LabelInstr           = "labelInstr"
	DeferInstr           = "deferInstr"
	TypeAliasInstr       = "typeAliasInstr"
)

type Instruction struct {
	Kind                InstructionKind  `json:"kind,omitempty"`
	Block               *Block           `json:"block,omitempty"`
	Assign              *Assign          `json:"assign,omitempty"`
	Local               *Local           `json:"local,omitempty"`
	Alloca              *Alloca          `json:"alloca,omitempty"`
	Next                *Next            `json:"next,omitempty"`
	Break               *Break           `json:"break,omitempty"`
	Return              *Return          `json:"return,omitempty"`
	Loop                *Loop            `json:"loop,omitempty"`
	Defer               *Defer           `json:"defer,omitempty"`
	Label               *Label           `json:"label,omitempty"`
	Jump                *Jump            `json:"jump,omitempty"`
	WhileLoop           *WhileLoop       `json:"whileLoop,omitempty"`
	ElseIfStatement     *ElseIfStatement `json:"elseIfStat,omitempty"`
	IfStatement         *IfStatement     `json:"ifStat,omitempty"`
	ExpressionStatement *Value           `json:"exprStat,omitempty"`
	TypeAliasStatement  *TypeAlias       `json:"typeAliasStat,omitempty"`
}

// TYPE ALIAS

// TypeAlias ...
// "type" iden "=" type ";"
type TypeAlias struct {
	Name front.Token
	Type *Type
}

func NewTypeAlias(name front.Token, typ *Type) *TypeAlias {
	return &TypeAlias{name, typ}
}

// BLOCK

type Block struct {
	ID         uint64         `json:"id,omitempty"`
	DeferStack []*Defer       `json:"deferStack,omitempty"`
	Instr      []*Instruction `json:"instr,omitempty"`
	Stab       *SymbolTable   `json:"stab,omitempty"`
	Return     *Instruction   `json:"return,omitempty"`
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

var blockIota uint64

func NewBlock() *Block {
	res := &Block{
		blockIota,
		[]*Defer{},
		[]*Instruction{},
		nil,
		nil,
	}
	blockIota++
	return res
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

type Label struct {
	Name front.Token
}

func NewLabel(name front.Token) *Label {
	return &Label{name}
}

type Jump struct {
	Location front.Token
}

func NewJump(loc front.Token) *Jump {
	return &Jump{loc}
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
	Owned   bool
	Val     *Value
}

// IsPublic will return the visiblity of the
// Local. If the first character of the name is
// uppercase, then the local is public
func (l *Local) IsPublic() bool {
	first := string(l.Name.Value[0])
	return strings.Compare(first, strings.ToUpper(first)) == 0
}

func (l *Local) SetValue(v *Value) {
	l.Val = v
}

func (l *Local) SetMutable(m bool) {
	l.Mutable = m
}

func NewLocal(name front.Token, typ *Type, owned bool) *Local {
	return &Local{name, typ, false, owned, nil}
}

// ALLOCA

type Alloca struct {
	Name    front.Token
	Type    *Type
	Mutable bool
	Owned   bool
	Val     *Value
}

// IsPublic will return the visiblity of the
// Alloca. If the first character of the name is
// uppercase, then the alloca is public
func (a *Alloca) IsPublic() bool {
	first := string(a.Name.Value[0])
	return strings.Compare(first, strings.ToUpper(first)) == 0
}

func (a *Alloca) SetValue(v *Value) {
	a.Val = v
}

func NewAlloca(name front.Token, mutable bool, owned bool, typ *Type) *Alloca {
	return &Alloca{name, typ, mutable, owned, nil}
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
	Cond   *Value             `json:"cond"`
	True   *Block             `json:"trueBlock"`
	ElseIf []*ElseIfStatement `json:"elseIf,omitempty"`
	Else   *Block             `json:"else,omitempty"`
}

func NewIfStatement(cond *Value, t *Block, elses []*ElseIfStatement, f *Block) *IfStatement {
	return &IfStatement{cond, t, elses, f}
}
