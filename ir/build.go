package ir

import (
	"fmt"
	"github.com/krug-lang/caasper/api"
	"github.com/krug-lang/caasper/front"
)

type builder struct {
	mod    *Module
	errors []api.CompilerError
}

func (b *builder) error(err api.CompilerError) {
	b.errors = append(b.errors, err)
}

func newBuilder(mod *Module) *builder {
	return &builder{mod, []api.CompilerError{}}
}

func (b *builder) buildUnresolvedType(u *front.UnresolvedTypeNode) *Type {
	if t, ok := PrimitiveType[u.Name]; ok {
		return t
	}
	return &Type{Kind: ReferenceKind, Reference: NewReferenceType(u.Name)}
}

func (b *builder) buildPointerType(p *front.PointerTypeNode) *PointerType {
	base := b.buildType(p.Base)
	return NewPointerType(base)
}

func (b *builder) buildArrayType(p *front.ArrayTypeNode) *Type {
	size := b.buildExpr(p.Size)

	base := b.buildType(p.Base)
	return &Type{
		Kind:      ArrayKind,
		ArrayType: NewArrayType(base, size),
	}
}

func (b *builder) buildStructureType(struc *front.StructureTypeNode) *Type {
	fields := newTypeDict()
	for _, sf := range struc.Fields {
		typ := b.buildType(sf.Type)
		fields.Add(NewLocal(sf.Name, typ, sf.Owned))
	}
	return &Type{
		Kind:      StructKind,
		Structure: NewStructure(struc.Name, fields),
	}
}

func (b *builder) buildTupleType(tup *front.TupleTypeNode) *Type {
	types := []*Type{}
	for _, typ := range tup.Types {
		types = append(types, b.buildType(typ))
	}
	return &Type{
		Kind:  TupleKind,
		Tuple: NewTupleType(types),
	}
}

func (b *builder) buildTypeExpr(te *front.TypeNode) *Type {
	switch te.Kind {
	case front.ArrayType:
		return b.buildArrayType(te.ArrayTypeNode)
	case front.UnresolvedType:
		return b.buildUnresolvedType(te.UnresolvedTypeNode)
	case front.TupleType:
		return b.buildTupleType(te.TupleTypeNode)
	case front.StructureType:
		return b.buildStructureType(te.StructureTypeNode)
	default:
		panic(fmt.Sprintf("unimplemented type_expr %s", te.Kind))
	}
}

func (b *builder) buildConstType(node *front.ConstantNode) *Type {
	resp := &Type{}

	switch node.Kind {
	case front.IntegerConstant:
		resp.IntegerType = NewIntegerType(32, true)
		resp.Kind = IntegerKind
	case front.FloatingConstant:
		resp.FloatingType = NewFloatingType(64)
		resp.Kind = FloatKind

	case front.VariableReference:
		resp.Kind = ReferenceKind

	default:
		panic(fmt.Sprintf("unimplemented buildConstType %s", node.Kind))
	}
	return resp
}

func (b *builder) buildType(node *front.ExpressionNode) *Type {
	resp := &Type{}

	switch node.Kind {
	case front.TypeExpression:
		resp = b.buildTypeExpr(node.TypeExpressionNode)

	case front.ConstantExpression:
		resp = b.buildConstType(node.ConstantNode)

	default:
		panic(fmt.Sprintf("unimplemented type %s", node.Kind))
	}

	return resp
}

func (b *builder) buildBinaryExpr(e *front.BinaryExpressionNode) *Value {
	lh := b.buildExpr(e.LHand)
	rh := b.buildExpr(e.RHand)
	return &Value{
		Kind:             BinaryExpressionValue,
		BinaryExpression: NewBinaryExpression(lh, e.Operator, rh),
	}
}

func (b *builder) buildUnaryExpr(u *front.UnaryExpressionNode) *Value {
	val := b.buildExpr(u.Value)
	return &Value{
		Kind:            UnaryExpressionValue,
		UnaryExpression: NewUnaryExpression(u.Operator, val),
	}
}

func (b *builder) buildGrouping(g *front.GroupingNode) *Value {
	val := b.buildExpr(g.Value)
	return &Value{
		Kind:     GroupingValue,
		Grouping: NewGrouping(val),
	}
}

func (b *builder) buildBuiltin(e *front.BuiltinExpressionNode) *Value {
	args := make([]*Value, len(e.Args))
	for i, arg := range e.Args {
		args[i] = b.buildExpr(arg)
	}
	return &Value{
		Kind:    BuiltinValue,
		Builtin: NewBuiltin(e.Name, NewIdentifier(e.Iden.Name), args),
	}
}

func (b *builder) buildCallExpression(e *front.CallExpressionNode) *Value {
	left := b.buildExpr(e.Left)
	var params []*Value
	for _, p := range e.Params {
		expr := b.buildExpr(p)
		params = append(params, expr)
	}
	return &Value{
		Kind: CallValue,
		Call: NewCall(left, params),
	}
}

func (b *builder) buildIndexExpression(i *front.IndexExpressionNode) *Value {
	left := b.buildExpr(i.Left)
	sub := b.buildExpr(i.Value)
	return &Value{
		Kind:  IndexValue,
		Index: NewIndex(left, sub),
	}
}

func (b *builder) buildPathExpression(p *front.PathExpressionNode) *Value {
	var values []*Value
	for _, e := range p.Values {
		val := b.buildExpr(e)

		if val.Kind == PathValue {
			path := val.Path
			for _, val := range path.Values {
				values = append(values, val)
			}
			continue
		}

		values = append(values, val)
	}

	return &Value{
		Kind: PathValue,
		Path: NewPath(values),
	}
}

func (b *builder) buildConst(e *front.ConstantNode) *Value {
	res := &Value{}

	switch e.Kind {
	case front.IntegerConstant:
		res.IntegerValue = NewIntegerValue(e.IntegerConstantNode.Value)
		res.Kind = IntegerValueValue
	case front.FloatingConstant:
		res.FloatingValue = NewFloatingValue(e.FloatingConstantNode.Value)
		res.Kind = FloatingValueValue
	case front.StringConstant:
		res.StringValue = NewStringValue(e.StringConstantNode.Value)
		res.Kind = StringValueValue
	case front.CharacterConstant:
		res.CharacterValue = NewCharacterValue(e.CharacterConstantNode.Value)
		res.Kind = CharacterValueValue
	case front.VariableReference:
		res.Kind = IdentifierValue
		res.Identifier = NewIdentifier(e.VariableReferenceNode.Name)
	default:
		panic(fmt.Sprintf("unimplemented %s", e.Kind))
	}
	return res
}

func (b *builder) buildLambda(expr *front.LambdaExpressionNode) *Value {
	// generate a new temp function
	// return a pointer to this new temp func?
	/*
		let x = fn(a int, b int) int {
			return a + b;
		};

		fn generatedFunc(a int, b int) int {
			return a + b;
		}

		let x = &generatedFunc;

		x();
	*/

	// hm.
	return &Value{}
}

func (b *builder) buildInitializerList(init *front.InitializerExpressionNode) *Value {
	vals := []*Value{}
	for _, expr := range init.Values {
		vals = append(vals, b.buildExpr(expr))
	}

	var iden *Identifier
	if init.Kind == front.InitStructure {
		iden = NewIdentifier(init.LHand)
	}

	return &Value{
		Kind: InitValue,
		Init: NewInit(init.Kind, iden, vals),
	}
}

func (b *builder) buildExpr(expr *front.ExpressionNode) *Value {
	switch expr.Kind {
	case front.ConstantExpression:
		return b.buildConst(expr.ConstantNode)

	case front.LambdaExpression:
		return b.buildLambda(expr.LambdaExpressionNode)

	case front.BinaryExpression:
		return b.buildBinaryExpr(expr.BinaryExpressionNode)

	case front.Grouping:
		return b.buildGrouping(expr.GroupingNode)

	case front.BuiltinExpression:
		return b.buildBuiltin(expr.BuiltinExpressionNode)

	case front.UnaryExpression:
		return b.buildUnaryExpr(expr.UnaryExpressionNode)

	case front.AssignStatement:
		return b.buildAssignStat(expr.AssignStatementNode)

	case front.CallExpression:
		return b.buildCallExpression(expr.CallExpressionNode)

	case front.PathExpression:
		{
			p := b.buildPathExpression(expr.PathExpressionNode)

			pat := p.Path
			last := pat.Values[len(pat.Values)-1]

			// rewrite so that we return a BinaryExpr(path, op, right)
			// versus
			// Path where the last value is a Binary Expr.
			if last.Kind == BinaryExpressionValue {
				bin := last.BinaryExpression

				// set the last value to be equal to the binary expr lhand
				pat.Values[len(pat.Values)-1] = bin.LHand
				bin.LHand = p

				return &Value{
					Kind:             BinaryExpressionValue,
					BinaryExpression: bin,
				}
			}

			return p
		}

	case front.IndexExpression:
		return b.buildIndexExpression(expr.IndexExpressionNode)

	case front.InitializerExpression:
		return b.buildInitializerList(expr.InitializerExpressionNode)

	default:
		panic(fmt.Sprintf("unhandled expr %s", expr.Kind))
	}

}

func (b *builder) buildLetStat(l *front.LetStatementNode) *Instruction {
	var val *Value
	if l.Value != nil {
		val = b.buildExpr(l.Value)
	}

	var typ *Type
	if l.Type != nil {
		typ = b.buildType(l.Type)
	} else if val != nil {
		// infer type from expr.
		// TODO annoying switch on kind for this
		// typ = val.InferredType()
	}

	local := NewLocal(l.Name, typ, l.Owned)
	local.SetValue(val)
	local.SetMutable(false)
	return &Instruction{
		Kind:  LocalInstr,
		Local: local,
	}
}

func (b *builder) buildMutStat(m *front.MutableStatementNode) *Instruction {
	var val *Value
	if m.Value != nil {
		val = b.buildExpr(m.Value)
	}

	var typ *Type
	if m.Type != nil {
		typ = b.buildType(m.Type)
	} else {
		if val == nil {
			panic("no expression to infer from")
		}
		// infer type from expr.
		// TODO annoying switch on kind for this
		// typ = val.InferredType()
	}

	local := NewLocal(m.Name, typ, m.Owned)
	local.SetValue(val)
	local.SetMutable(true)
	return &Instruction{
		Kind:  LocalInstr,
		Local: local,
	}
}

func (b *builder) buildReturnStat(ret *front.ReturnStatementNode) *Instruction {
	var val *Value
	if ret.Value != nil {
		val = b.buildExpr(ret.Value)
	}
	res := NewReturn(val)
	return &Instruction{
		Kind:   ReturnInstr,
		Return: res,
	}
}

func (b *builder) buildWhileLoopStat(while *front.WhileLoopNode) *Instruction {
	cond := b.buildExpr(while.Cond)
	var post *Value
	if while.Post != nil {
		post = b.buildExpr(while.Post)
	}
	body := b.buildBlock(while.Block)
	res := NewWhileLoop(cond, post, body)
	return &Instruction{
		Kind:      WhileLoopInstr,
		WhileLoop: res,
	}
}

func (b *builder) buildLoopStat(loop *front.LoopNode) *Instruction {
	body := b.buildBlock(loop.Block)
	res := NewLoop(body)
	return &Instruction{
		Kind: LoopInstr,
		Loop: res,
	}
}

func (b *builder) buildBlock(block *front.BlockNode) *Block {
	res := NewBlock()

	for _, stat := range block.Statements {
		st := b.buildStat(stat)
		if st == nil {
			continue
		}

		switch st.Kind {

		case DeferInstr:
			// do we add this to the instructions list?
			res.PushDefer(st.Defer)

		case ReturnInstr:
			res.SetReturn(st)
			fallthrough
		default:
			res.AddInstr(st)
		}
	}
	return res
}

func (b *builder) buildIfStat(iff *front.IfNode) *Instruction {
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

	return &Instruction{
		Kind:        IfStatementInstr,
		IfStatement: NewIfStatement(cond, t, elses, f),
	}
}

func (b *builder) buildAssignStat(a *front.AssignStatementNode) *Value {
	// FIXME!
	res := NewAssign(b.buildExpr(a.LHand), a.Op, b.buildExpr(a.RHand))
	return &Value{
		Kind:   AssignValue,
		Assign: res,
	}
}

func (b *builder) buildDefer(def *front.DeferNode) *Instruction {
	var stat *Instruction
	var block *Block

	if def.Statement != nil {
		stat = b.buildStat(def.Statement)
	}
	if def.Block != nil {
		block = b.buildBlock(def.Block)
	}

	return &Instruction{
		Kind:  DeferInstr,
		Defer: NewDefer(stat, block),
	}
}

func (b *builder) buildStat(stat *front.ParseTreeNode) *Instruction {
	switch stat.Kind {
	case front.BlockStatement:
		body := b.buildBlock(stat.BlockNode)
		return &Instruction{
			Kind:  BlockInstr,
			Block: body,
		}

	case front.LetStatement:
		return b.buildLetStat(stat.LetStatementNode)
	case front.MutableStatement:
		return b.buildMutStat(stat.MutableStatementNode)

	case front.ReturnStatement:
		return b.buildReturnStat(stat.ReturnStatementNode)

	case front.BreakStatement:
		return &Instruction{Kind: BreakInstr, Break: NewBreak()}
	case front.NextStatement:
		return &Instruction{Kind: NextInstr, Next: NewNext()}

	case front.LoopStatement:
		return b.buildLoopStat(stat.LoopNode)
	case front.WhileLoopStatement:
		return b.buildWhileLoopStat(stat.WhileLoopNode)

	case front.IfStatement:
		return b.buildIfStat(stat.IfNode)

	case front.DeferStatement:
		return b.buildDefer(stat.DeferNode)

	case front.ExpressionStatement:
		return &Instruction{
			Kind:                ExpressionInstr,
			ExpressionStatement: b.buildExpr(stat.ExpressionStatementNode),
		}

	case front.TypeAliasStatement:
		typ := b.buildType(stat.TypeAliasNode.Type)
		return &Instruction{
			Kind:               TypeAliasInstr,
			TypeAliasStatement: NewTypeAlias(stat.TypeAliasNode.Name, typ),
		}

	case front.JumpStatement:
		return &Instruction{
			Kind: JumpInstr,
			Jump: NewJump(stat.JumpNode.Location),
		}

	case front.LabelStatement:
		return &Instruction{
			Kind:  LabelInstr,
			Label: NewLabel(stat.LabelNode.LabelName),
		}

	default:
		panic(fmt.Sprintf("unimplemented stat! %s", stat.Kind))
	}
}

func (b *builder) buildFunc(node *front.FunctionDeclaration) *Function {
	params := newTypeDict()
	for _, p := range node.Arguments {
		param := NewLocal(p.Name, b.buildType(p.Type), p.Owned)
		param.SetMutable(p.Mutable)
		params.Add(param)
	}

	ret := Void
	if node.ReturnType != nil {
		ret = b.buildType(node.ReturnType)
	}

	fn := NewFunction(node.Name, params, ret)
	fn.Body = b.buildBlock(node.Body)
	return fn
}

func (b *builder) buildTypeAlias(nt *front.TypeAliasNode) *Instruction {
	typ := b.buildType(nt.Type)
	return &Instruction{
		Kind:               TypeAliasInstr,
		TypeAliasStatement: NewTypeAlias(nt.Name, typ),
	}
}

func (b *builder) introduceNamedTypes(nodes []*front.ParseTreeNode) {
	for _, node := range nodes {
		if node.Kind != front.TypeAliasStatement {
			continue
		}
		alias := b.buildTypeAlias(node.TypeAliasNode)
		b.mod.Global.AddInstr(alias)
	}
}

func (b *builder) buildFunctions(nodes []*front.ParseTreeNode) {
	for _, node := range nodes {
		if node.Kind != front.FunctionDeclStatement {
			continue
		}
		fn := b.buildFunc(node.FunctionDeclaration)
		b.mod.RegisterFunction(fn)
	}
}

// buildTree takes the given set of nodes, and builds a module
// from them.
func (b *builder) buildTree(m *Module, nodes []*front.ParseTreeNode) {
	b.introduceNamedTypes(nodes)
	b.buildFunctions(nodes)
}

func Build(trees [][]*front.ParseTreeNode) (*Module, []api.CompilerError) {
	module := NewModule("main")

	b := newBuilder(module)
	for _, tree := range trees {
		fmt.Println(tree)
		b.buildTree(module, tree)
	}
	return module, b.errors
}
