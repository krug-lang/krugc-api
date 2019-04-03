package ir

import (
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/front"
)

// Build takes the given []ParseTree's
// and builds a SINGLE ir module from them.
func Build(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var trees [][]*front.ParseTreeNode

	if err := jsoniter.Unmarshal([]byte(krugReq.Data), &trees); err != nil {
		panic(err)
	}

	irModule, errors := build(trees)

	jsonIrModule, err := jsoniter.MarshalIndent(irModule, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonIrModule),
		Errors: errors,
	}
	c.JSON(200, &resp)
}

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
	// TODO should be constant expr.
	size := b.buildExpr(p.Size)

	base := b.buildType(p.Base)
	return &Type{
		Kind:      ArrayKind,
		ArrayType: NewArrayType(base, size),
	}
}

func (b *builder) buildType(node *front.TypeNode) *Type {
	var resp *Type

	switch node.Kind {
	case front.UnresolvedType:
		return b.buildUnresolvedType(node.UnresolvedTypeNode)

	case front.PointerType:
		resp.Pointer = b.buildPointerType(node.PointerTypeNode)
		resp.Kind = PointerKind

	case front.ArrayType:
		return b.buildArrayType(node.ArrayTypeNode)

	default:
		panic(fmt.Sprintf("unimplemented type %s", reflect.TypeOf(node)))
	}

	return resp
}

// here we build a structure, give it an EMPTY type dictionary
// of the fields.
func (b *builder) declareStructure(node *front.StructureDeclaration) *Structure {
	fields := newTypeDict()
	return NewStructure(node.Name, fields)
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
	return &Value{
		Kind:    BuiltinValue,
		Builtin: NewBuiltin(e.Name, b.buildType(e.Type)),
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

		// stupid hack to flatten the path expressions
		// TODO(felix): fix the parser for this.
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
	default:
		panic(fmt.Sprintf("unimplemented %s", e.Kind))
	}
	return res
}

func (b *builder) buildExpr(expr *front.ExpressionNode) *Value {
	switch expr.Kind {
	case front.ConstantExpression:
		return b.buildConst(expr.ConstantNode)

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
		return b.buildPathExpression(expr.PathExpressionNode)

	case front.IndexExpression:
		return b.buildIndexExpression(expr.IndexExpressionNode)

	default:
		panic(fmt.Sprintf("unhandled expr %s", reflect.TypeOf(expr)))
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

	local := NewLocal(l.Name, typ)
	local.SetValue(val)
	local.SetMutable(true)
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

	local := NewLocal(m.Name, typ)
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
		if st := b.buildStat(stat); st != nil {
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

	case front.ExpressionStatement:
		return &Instruction{
			Kind:                ExpressionInstr,
			ExpressionStatement: b.buildExpr(stat.ExpressionStatementNode),
		}

	default:
		panic(fmt.Sprintf("unimplemented stat! %s", reflect.TypeOf(stat)))
	}
}

func (b *builder) buildFunc(node *front.FunctionDeclaration) *Function {
	params := newTypeDict()
	for _, p := range node.Arguments {
		params.Set(p.Name, b.buildType(p.Type))
	}

	ret := Void
	if node.ReturnType != nil {
		ret = b.buildType(node.ReturnType)
	}

	fn := NewFunction(node.Name, params, ret)
	fn.Body = b.buildBlock(node.Body)
	return fn
}

func (b *builder) buildTree(m *Module, nodes []*front.ParseTreeNode) {
	structureNodes := []*front.StructureDeclaration{}
	implNodes := []*front.ImplDeclaration{}

	// TODO it might be faster to do one loop and append
	// all the types into arrays then process the arrays.

	// declare all structures, impls, exist.
	for _, n := range nodes {
		switch n.Kind {
		case front.StructureDeclStatement:
			struc := n.StructureDeclaration
			structureNodes = append(structureNodes, struc)
			m.RegisterStructure(b.declareStructure(struc))

		case front.ImplDeclStatement:
			impl := n.ImplDeclaration
			implNodes = append(implNodes, impl)
			if m.RegisterImpl(NewImpl(impl.Name)) {
				b.error(api.CompilerError{
					Title: fmt.Sprintf("Duplicate implementation for '%s'", impl.Name),
					Desc:  "...",
				})
			}
		}
	}

	// go through structures again
	// this time process the fields
	for _, sn := range structureNodes {
		structure, ok := m.GetStructure(sn.Name.Value)
		if !ok {
			panic(fmt.Sprintf("couldn't find structure %s", sn.Name))
		}

		for _, tn := range sn.Fields {
			typ := b.buildType(tn.Type)
			if typ == nil {
				panic("couldn't build type when setting structure field")
			}
			structure.Fields.Set(tn.Name, typ)
		}
	}

	for _, in := range implNodes {
		impl, ok := m.GetImpl(in.Name.Value)
		if !ok {
			panic(fmt.Sprintf("couldn't find impl %s", in.Name))
		}

		for _, fn := range in.Functions {
			builtFunc := b.buildFunc(fn)
			ok := impl.RegisterMethod(builtFunc)
			if !ok {
				b.error(api.NewSymbolError(fn.Name.Value, fn.Name.Span...))
			}
		}
	}

	// TODO do the same as above but for functions
	// i.e. declare then define.

	// then we do all the functions
	for _, n := range nodes {
		switch n.Kind {
		case front.FunctionDeclStatement:
			f := b.buildFunc(n.FunctionDeclaration)
			m.RegisterFunction(f)
		}
	}
}

func build(trees [][]*front.ParseTreeNode) (*Module, []api.CompilerError) {
	module := NewModule("main")

	b := newBuilder(module)
	for _, tree := range trees {
		fmt.Println(tree)
		b.buildTree(module, tree)
	}
	return module, b.errors
}
