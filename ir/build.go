package ir

import (
	"fmt"
	"reflect"

	"github.com/krug-lang/krugc-api/front"
)

var typeMap = map[string]Type{
	"u8":  Uint8,
	"u16": Uint16,
	"u32": Uint32,
	"u64": Uint64,

	"s8":  Int8,
	"s16": Int16,
	"s32": Int32,
	"s64": Int64,

	"void": Void,
	"bool": Bool,

	"f32": Float32,
	"f64": Float64,

	"rune": Int32,

	"int":  Int32,
	"uint": Uint32,
}

type builder struct {
	mod *Module
}

func newBuilder(mod *Module) *builder {
	return &builder{mod}
}

func (b *builder) buildUnresolvedType(u *front.UnresolvedType) Type {
	if t, ok := typeMap[u.Name]; ok {
		return t
	}

	if _, ok := b.mod.GetStructure(u.Name); ok {
		// FIXME: should we store the type in this reference?
		// probably should.
		return NewReferenceType(u.Name)
	}

	panic(fmt.Sprintf("not sure where the type %s comes from", u.Name))
}

func (b *builder) buildPointerType(p *front.PointerType) Type {
	base := b.buildType(p.Base)
	return NewPointerType(base)
}

func (b *builder) buildType(n front.TypeNode) Type {
	switch node := n.(type) {
	case *front.UnresolvedType:
		return b.buildUnresolvedType(node)

	case *front.PointerType:
		return b.buildPointerType(node)

	default:
		panic(fmt.Sprintf("unimplemented type %s", reflect.TypeOf(node)))
	}
}

func (b *builder) buildStructure(node *front.StructureDeclaration) *Structure {
	fields := map[string]Type{}
	return NewStructure(node.Name, fields)
}

func (b *builder) buildBinaryExpr(e *front.BinaryExpression) *BinaryExpression {
	lh := b.buildExpr(e.LHand)
	rh := b.buildExpr(e.RHand)
	return NewBinaryExpression(lh, e.Operator, rh)
}

func (b *builder) buildUnaryExpr(u *front.UnaryExpression) *UnaryExpression {
	val := b.buildExpr(u.Value)
	return NewUnaryExpression(u.Operator, val)
}

func (b *builder) buildGrouping(g *front.Grouping) *Grouping {
	val := b.buildExpr(g.Value)
	return NewGrouping(val)
}

func (b *builder) buildBuiltin(e *front.BuiltinExpression) Value {
	return NewBuiltin(e.Name, b.buildType(e.Type))
}

func (b *builder) buildExpr(e front.ExpressionNode) Value {
	switch expr := e.(type) {
	case *front.IntegerConstant:
		return NewIntegerValue(expr.Value)

	case *front.BinaryExpression:
		return b.buildBinaryExpr(expr)

	case *front.VariableReference:
		return NewIdentifier(expr.Name)

	case *front.Grouping:
		return b.buildGrouping(expr)

	case *front.BuiltinExpression:
		return b.buildBuiltin(expr)

	case *front.UnaryExpression:
		return b.buildUnaryExpr(expr)

	case *front.AssignStatement:
		return b.buildAssignStat(expr)

	default:
		panic(fmt.Sprintf("unhandled expr %s", reflect.TypeOf(expr)))
	}

}

func (b *builder) buildLetStat(l *front.LetStatement) Instruction {
	typ := b.buildType(l.Type)
	local := NewLocal(l.Name, typ)
	if l.Value != nil {
		local.SetValue(b.buildExpr(l.Value))
	}
	return local
}

func (b *builder) buildMutStat(m *front.MutableStatement) Instruction {
	typ := b.buildType(m.Type)
	local := NewLocal(m.Name, typ)
	if m.Value != nil {
		local.SetValue(b.buildExpr(m.Value))
	}
	local.SetMutable(true)
	return local
}

func (b *builder) buildReturnStat(ret *front.ReturnStatement) Instruction {
	var val Value
	if ret.Value != nil {
		val = b.buildExpr(ret.Value)
	}
	return NewReturn(val)
}

func (b *builder) buildWhileLoopStat(while *front.WhileLoopStatement) Instruction {
	cond := b.buildExpr(while.Cond)
	var post Value
	if while.Post != nil {
		post = b.buildExpr(while.Post)
	}
	body := b.buildBlock(while.Block)
	return NewWhileLoop(cond, post, body)
}

func (b *builder) buildLoopStat(loop *front.LoopStatement) Instruction {
	body := b.buildBlock(loop.Block)
	return NewLoop(body)
}

func (b *builder) buildBlock(stats []front.StatementNode) *Block {
	block := NewBlock()
	for _, stat := range stats {
		if st := b.buildStat(stat); st != nil {
			block.AddInstr(st)
		}
	}
	return block
}

func (b *builder) buildIfStat(iff *front.IfStatement) Instruction {
	cond := b.buildExpr(iff.Cond)
	t := b.buildBlock(iff.Block)

	elses := []*ElseIfStatement{}
	for _, elif := range iff.ElseIfs {
		cond := b.buildExpr(elif.Cond)
		block := b.buildBlock(elif.Block)
		elses = append(elses, NewElseIfStatement(cond, block))
	}

	var f *Block
	if iff.Else != nil {
		f = b.buildBlock(iff.Else)
	}

	return NewIfStatement(cond, t, elses, f)
}

func (b *builder) buildAssignStat(a *front.AssignStatement) Instruction {
	// FIXME!
	return NewAssign(b.buildExpr(a.LHand), a.Op, b.buildExpr(a.RHand))
}

func (b *builder) buildStat(s front.StatementNode) Instruction {
	switch stat := s.(type) {
	case *front.LetStatement:
		return b.buildLetStat(stat)
	case *front.MutableStatement:
		return b.buildMutStat(stat)
	case *front.ReturnStatement:
		return b.buildReturnStat(stat)
	case *front.LoopStatement:
		return b.buildLoopStat(stat)
	case *front.WhileLoopStatement:
		return b.buildWhileLoopStat(stat)
	case *front.IfStatement:
		return b.buildIfStat(stat)
	case *front.AssignStatement:
		return b.buildAssignStat(stat)
	default:
		panic(fmt.Sprintf("unimplemented stat! %s", reflect.TypeOf(stat)))
	}
}

func (b *builder) buildFunc(node *front.FunctionDeclaration) *Function {
	params := map[string]Type{}
	for n, p := range node.Arguments {
		params[n] = b.buildType(p)
	}

	var ret Type = Void
	if node.ReturnType != nil {
		ret = b.buildType(node.ReturnType)
	}

	fn := NewFunction(node.Name, params, ret)
	fn.Body = b.buildBlock(node.Body)
	return fn
}

func (b *builder) buildTree(m *Module, tree front.ParseTree) {
	structureNodes := []*front.StructureDeclaration{}

	// declare all structures
	for _, n := range tree.Nodes {
		switch node := n.(type) {
		case *front.StructureDeclaration:
			structureNodes = append(structureNodes, node)
			m.RegisterStructure(b.buildStructure(node))
		}
	}

	// go through structures again
	// this time process the fields
	for _, sn := range structureNodes {
		structure, ok := m.GetStructure(sn.Name)
		if !ok {
			panic("couldn't find structure!")
		}

		for n, tn := range sn.Fields {
			structure.Fields[n] = b.buildType(tn)
		}
	}

	// TODO do the same as above but for functions
	// i.e. declare then define.

	// then we do all the functions
	for _, n := range tree.Nodes {
		switch node := n.(type) {
		case *front.FunctionDeclaration:
			f := b.buildFunc(node)
			m.RegisterFunction(f)
		}
	}
}

func build(trees []front.ParseTree) *Module {
	module := NewModule("main")

	b := newBuilder(module)
	for _, tree := range trees {
		b.buildTree(module, tree)
	}
	return module
}
