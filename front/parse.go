package front

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/hugobrains/caasper/api"
)

// keywords
const (
	fn    string = "fn"
	let          = "let"
	mut          = "mut"
	brk          = "break"
	ret          = "return"
	next         = "next"
	trait        = "trait"
	struc        = "struct"
	impl         = "impl"
)

var BadToken Token

type parser struct {
	toks   []Token
	pos    int
	errors []api.CompilerError
}

func (p *parser) error(e api.CompilerError) {
	p.errors = append(p.errors, e)
}

func (p *parser) peek(offs int) (tok Token) {
	tok = p.toks[p.pos+offs]
	return tok
}

func (p *parser) next() (tok Token) {
	tok = p.toks[p.pos]
	return tok
}

func (p *parser) expect(val string) (tok Token) {
	start := p.pos
	if p.hasNext() {
		if tok = p.consume(); tok.Matches(val) {
			return tok
		}
	}

	p.error(api.NewUnexpectedToken(val, p.next().Value, start, p.pos))
	return BadToken
}

func (p *parser) expectKind(kind TokenType) (tok Token) {
	start := p.pos

	if tok = p.consume(); tok.Kind == kind {
		return tok
	}

	p.error(api.NewUnexpectedToken(string(kind), p.next().Value, start, p.pos))
	return BadToken
}

func (p *parser) consume() (tok Token) {
	tok = p.toks[p.pos]
	p.pos++
	return tok
}

func (p *parser) hasNext() bool {
	return p.pos < len(p.toks)
}

func (p *parser) parsePointerType() *TypeNode {
	start := p.pos
	p.expect("*")
	base := p.parseType()
	if base == nil {
		p.error(api.NewParseError("type after pointer", start, p.pos))
		return nil
	}

	return &TypeNode{
		Kind: PointerType,
		PointerTypeNode: &PointerTypeNode{
			Base: base,
		},
	}
}

func (p *parser) parseArrayType() *TypeNode {
	start := p.pos

	p.expect("[")
	base := p.parseType()
	if base == nil {
		p.error(api.NewParseError("array type", start, p.pos))
	}

	p.expect(";")

	size := p.parseExpression()
	if size == nil {
		p.error(api.NewParseError("array length constant", start, p.pos))
	}

	p.expect("]")
	return &TypeNode{
		Kind: ArrayType,
		ArrayTypeNode: &ArrayTypeNode{
			Base: base,
			Size: size,
		},
	}
}

func (p *parser) parseUnresolvedType() *TypeNode {
	name := p.expectKind(Identifier)
	return &TypeNode{
		Kind: UnresolvedType,
		UnresolvedTypeNode: &UnresolvedTypeNode{
			Name: name.Value,
		},
	}
}

func (p *parser) parseType() *TypeNode {
	start := p.pos
	switch curr := p.next(); {
	case curr.Matches("^"):
		return p.parsePointerType()
	case curr.Matches("["):
		return p.parseArrayType()
	case curr.Kind == Identifier:
		return p.parseUnresolvedType()
	default:
		p.error(api.NewUnimplementedError("type", start, p.pos))
		return nil
	}
}

func (p *parser) parseStructureDeclaration() *StructureDeclaration {
	start := p.pos

	p.expect("struct")
	name := p.consume()

	fields := []*NamedType{}

	p.expect("{")
	for p.hasNext() {
		if p.next().Matches("}") {
			break
		}

		name := p.expectKind(Identifier)

		typ := p.parseType()
		if typ == nil {
			p.error(api.NewParseError("type", start, p.pos))
		}

		// FIXME
		fields = append(fields, &NamedType{name, typ})

		// trailing commas are enforced.
		p.expect(",")
	}
	p.expect("}")

	return &StructureDeclaration{
		Name:   name,
		Fields: fields,
	}
}

func (p *parser) parseFunctionPrototypeDeclaration() *FunctionPrototypeDeclaration {
	start := p.pos

	p.expect(fn)
	name := p.expectKind(Identifier)

	args := []*NamedType{}

	p.expect("(")
	for idx := 0; p.hasNext(); idx++ {
		if p.next().Matches(")") {
			break
		}

		// no trailing commas allowed here.
		if idx != 0 {
			p.expect(",")
		}

		name := p.expectKind(Identifier)
		typ := p.parseType()
		if typ == nil {
			p.error(api.NewParseError("type after pointer", start, p.pos))
		}

		args = append(args, &NamedType{name, typ})
	}
	p.expect(")")

	typ := p.parseType()

	return &FunctionPrototypeDeclaration{
		Name:      name,
		Arguments: args,

		// could be nil!
		ReturnType: typ,
	}
}

// mut x [ type ] [ = val ]
func (p *parser) parseMut() *ParseTreeNode {
	start := p.pos

	p.expect(mut)
	name := p.expectKind(Identifier)

	var typ *TypeNode
	if !p.next().Matches("=") {
		typ = p.parseType()
		if typ == nil {
			p.error(api.NewParseError("type after assignment", start, p.pos))
		}
	}

	var val *ExpressionNode
	if p.next().Matches("=") {
		p.expect("=")

		val = p.parseExpression()
		if val == nil {
			p.error(api.NewParseError("assignment", start, p.pos))
		}
	}

	if val == nil && typ == nil {
		p.error(api.NewParseError("value or type in mut statement", start, p.pos))
	}

	return &ParseTreeNode{
		Kind: MutableStatement,
		MutableStatementNode: &MutableStatementNode{
			Name:  name,
			Type:  typ,
			Value: val,
		},
	}
}

// let is a constant variable.
func (p *parser) parseLet() *ParseTreeNode {
	start := p.pos

	p.expect(let)
	name := p.expectKind(Identifier)

	var typ *TypeNode
	if !p.next().Matches("=") {
		typ = p.parseType()
		if typ == nil {
			p.error(api.NewParseError("type or assignment", start, p.pos))
		}
	}

	var value *ExpressionNode
	if p.next().Matches("=") {
		p.expect("=")
		value = p.parseExpression()
		if value == nil {
			p.error(api.NewParseError("expression in let statement", start, p.pos))
		}
	}

	return &ParseTreeNode{
		Kind: LetStatement,
		LetStatementNode: &LetStatementNode{
			Name:  name,
			Type:  typ,
			Value: value,
		},
	}
}

func (p *parser) parseReturn() *ParseTreeNode {
	if !p.next().Matches("return") {
		return nil
	}
	start := p.pos
	p.expect("return")

	var res *ExpressionNode
	if !p.next().Matches(";") {
		res = p.parseExpression()
		if res == nil {
			p.error(api.NewParseError("semi-colon or expression", start, p.pos))
		}
	}

	return &ParseTreeNode{
		Kind: ReturnStatement,
		ReturnStatementNode: &ReturnStatementNode{
			Value: res,
		},
	}
}

func (p *parser) parseNext() *ParseTreeNode {
	p.expect("next")
	return &ParseTreeNode{
		Kind: NextStatement,
	}
}

func (p *parser) parseBreak() *ParseTreeNode {
	p.expect("break")
	return &ParseTreeNode{
		Kind: BreakStatement,
	}
}

func (p *parser) parseSemicolonStatement() *ParseTreeNode {
	switch curr := p.next(); {
	case curr.Matches(mut):
		return p.parseMut()
	case curr.Matches(let):
		return p.parseLet()
	case curr.Matches(ret):
		return p.parseReturn()
	case curr.Matches(next):
		return p.parseNext()
	case curr.Matches(brk):
		return p.parseBreak()
	}

	return p.parseExpressionStatement()
}

func (p *parser) parseStatBlock() *BlockNode {
	if !p.next().Matches("{") {
		return nil
	}

	stats := []*ParseTreeNode{}
	p.expect("{")
	for p.hasNext() {
		if p.next().Matches("}") {
			break
		}

		if stat := p.parseStatement(); stat != nil {
			stats = append(stats, stat)
		}
	}
	p.expect("}")

	return &BlockNode{
		Statements: stats,
	}
}

func (p *parser) parseIfElseChain() *ParseTreeNode {
	if !p.next().Matches("if") {
		return nil
	}

	start := p.pos

	// the first if.
	p.expect("if")
	expr := p.parseExpression()
	if expr == nil {
		p.error(api.NewParseError("condition", start, p.pos))
	}

	block := p.parseStatBlock()
	if block == nil {
		p.error(api.NewParseError("block after condition", start, p.pos))
	}

	var elseBlock *BlockNode
	elses := []*ElseIfNode{}

	for p.hasNext() && p.next().Matches("else") {
		// else if
		if p.peek(1).Matches("if") {
			p.expect("else")
			p.expect("if")
			cond := p.parseExpression()
			if cond == nil {
				p.error(api.NewParseError("condition in else if", start, p.pos))
			}

			body := p.parseStatBlock()
			if body == nil {
				p.error(api.NewParseError("block after else if", start, p.pos))
			}
			elses = append(elses, &ElseIfNode{cond, body})
		} else {
			// TODO we could easily check if else has been set before
			// or should we do this in a later pass during sema analyis?

			p.expect("else")
			elseBlock = p.parseStatBlock()
			if elseBlock == nil {
				p.error(api.NewParseError("block after else", start, p.pos))
			}
		}
	}

	return &ParseTreeNode{
		Kind: IfStatement,
		IfNode: &IfNode{
			Cond:    expr,
			Block:   block,
			Else:    elseBlock,
			ElseIfs: elses,
		},
	}
}

func (p *parser) parseWhileLoop() *ParseTreeNode {
	if !p.next().Matches("while") {
		return nil
	}
	start := p.pos

	p.expect("while")
	val := p.parseExpression()
	if val == nil {
		p.error(api.NewParseError("condition after while", start, p.pos))
	}

	var post *ExpressionNode
	if p.next().Matches(";") {
		p.expect(";")
		post = p.parseExpression()
		if post == nil {
			p.error(api.NewParseError("step expression in while loop", start, p.pos))
		}
	}

	if block := p.parseStatBlock(); block != nil {
		return &ParseTreeNode{
			Kind: WhileLoopStatement,
			WhileLoopNode: &WhileLoopNode{
				Cond:  val,
				Post:  post,
				Block: block,
			},
		}
	}

	return nil
}

func (p *parser) parseLoop() *ParseTreeNode {
	if !p.next().Matches("loop") {
		return nil
	}
	p.expect("loop")
	if block := p.parseStatBlock(); block != nil {
		return &ParseTreeNode{
			Kind: LoopStatement,
			LoopNode: &LoopNode{
				Block: block,
			},
		}
	}
	return nil
}

func (p *parser) parseStatement() *ParseTreeNode {
	switch curr := p.next(); {
	case curr.Matches("if"):
		return p.parseIfElseChain()
	case curr.Matches("loop"):
		return p.parseLoop()
	case curr.Matches("while"):
		return p.parseWhileLoop()
	case curr.Matches("{"):
		return &ParseTreeNode{
			Kind:      BlockStatement,
			BlockNode: p.parseStatBlock(),
		}
	}

	stat := p.parseSemicolonStatement()
	if stat != nil {
		p.expect(";")
	}
	return stat
}

func (p *parser) parseFunctionDeclaration() *FunctionDeclaration {
	proto := p.parseFunctionPrototypeDeclaration()

	body := p.parseStatBlock()
	return &FunctionDeclaration{proto, body}
}

func (p *parser) parseImplDeclaration() *ImplDeclaration {
	p.expect("impl")
	name := p.expectKind(Identifier)

	functions := []*FunctionDeclaration{}

	p.expect("{")
	for p.hasNext() {
		if p.next().Matches("}") {
			break
		}

		// NOTE: we dont care if impls are empty.
		if fn := p.parseFunctionDeclaration(); fn != nil {
			functions = append(functions, fn)
		}
	}
	p.expect("}")

	return &ImplDeclaration{
		name, functions,
	}
}

func (p *parser) parseTraitDeclaration() *TraitDeclaration {
	p.expect("trait")

	name := p.expectKind(Identifier)

	members := []*FunctionPrototypeDeclaration{}

	p.expect("{")
	for p.hasNext() {
		if p.next().Matches("}") {
			break
		}

		// we only parse prototypes here, not
		// function bodies.

		pt := p.parseFunctionPrototypeDeclaration()
		if pt == nil {
			break
		}
		members = append(members, pt)

		// must have semi-colon after each prototype.
		p.expect(";")
	}
	p.expect("}")

	return &TraitDeclaration{name, members}
}

func (p *parser) parseUnaryExpr() *ExpressionNode {
	// TODO other unary ops.
	if !p.hasNext() || !p.next().Matches("-", "!", "+", "@") {
		return nil
	}

	start := p.pos

	op := p.consume()
	right := p.parseLeft()
	if right == nil {
		p.error(api.NewParseError("unary expression", start, p.pos))
	}

	return &ExpressionNode{
		Kind:                UnaryExpression,
		UnaryExpressionNode: &UnaryExpressionNode{op.Value, right},
	}
}

func (p *parser) parseOperand() *ExpressionNode {
	if !p.hasNext() {
		return nil
	}

	start := p.pos
	curr := p.next()

	if curr.Matches("(") {
		p.expect("(")
		expr := p.parseExpression()
		p.expect(")")
		return &ExpressionNode{
			Kind:         Grouping,
			GroupingNode: &GroupingNode{expr},
		}
	}

	switch curr := p.consume(); curr.Kind {
	case Number:
		// no dot means it's a whole number.
		if strings.Index(curr.Value, ".") == -1 {
			bigint := new(big.Int)
			bigint.SetString(curr.Value, 10)

			return &ExpressionNode{
				Kind: ConstantExpression,
				ConstantNode: &ConstantNode{
					Kind:                IntegerConstant,
					IntegerConstantNode: &IntegerConstantNode{bigint},
				},
			}
		}

		val, err := strconv.ParseFloat(curr.Value, 64)
		if err != nil {
			panic(err)
		}

		return &ExpressionNode{
			Kind: ConstantExpression,
			ConstantNode: &ConstantNode{
				Kind:                 FloatingConstant,
				FloatingConstantNode: &FloatingConstantNode{val},
			},
		}
	case Identifier:
		return &ExpressionNode{
			Kind: ConstantExpression,
			ConstantNode: &ConstantNode{
				Kind:                  VariableReference,
				VariableReferenceNode: &VariableReferenceNode{curr},
			},
		}

	case String:
		return &ExpressionNode{
			Kind: ConstantExpression,
			ConstantNode: &ConstantNode{
				Kind:               StringConstant,
				StringConstantNode: &StringConstantNode{curr.Value},
			},
		}

	case EndOfFile:
		return nil

	default:
		p.error(api.NewUnimplementedError(curr.Value, start, p.pos))
		return nil
	}
}

func (p *parser) parseBuiltin() *ExpressionNode {
	start := p.pos
	builtin := p.expectKind(Identifier)
	p.expect("!")

	parens := false
	if p.next().Matches("(") {
		parens = true
		p.consume()
	}

	typ := p.parseType()
	if typ == nil {
		p.error(api.NewParseError("type in builtin", start, p.pos))
	}

	if parens {
		p.expect(")")
	}

	return &ExpressionNode{
		Kind:                  BuiltinExpression,
		BuiltinExpressionNode: &BuiltinExpressionNode{builtin.Value, typ},
	}
}

func (p *parser) parseCall(left *ExpressionNode) *ExpressionNode {
	start := p.pos

	var params []*ExpressionNode

	p.expect("(")
	for idx := 0; p.hasNext() && !p.next().Matches(")"); idx++ {
		if idx != 0 {
			p.expect(",")
		}

		val := p.parseExpression()
		if val == nil {
			p.error(api.NewParseError("parameter in call expression", start, p.pos))
		}
		params = append(params, val)
	}
	p.expect(")")

	return &ExpressionNode{
		Kind: CallExpression,
		CallExpressionNode: &CallExpressionNode{
			left, params,
		},
	}
}

func (p *parser) parseIndex(left *ExpressionNode) *ExpressionNode {
	start := p.pos
	p.expect("[")
	val := p.parseExpression()
	if val == nil {
		p.error(api.NewParseError("expression in array index", start, p.pos))
	}
	p.expect("]")
	return &ExpressionNode{
		Kind: IndexExpression,
		IndexExpressionNode: &IndexExpressionNode{
			left, val,
		},
	}
}

func (p *parser) parseLambda() *ExpressionNode {
	proto := p.parseFunctionPrototypeDeclaration()
	body := p.parseStatBlock()
	return &ExpressionNode{
		Kind: LambdaExpression,
		LambdaExpressionNode: &LambdaExpressionNode{
			proto, body,
		},
	}
}

func (p *parser) parsePrimaryExpr() *ExpressionNode {
	if !p.hasNext() {
		return nil
	}

	if p.next().Matches(fn) {
		return p.parseLambda()
	}

	// TODO unary ops.
	if p.next().Matches("!", "@", "+", "-") {
		return p.parseUnaryExpr()
	}

	// builtins.
	switch curr := p.next(); {
	case curr.Matches("make", "sizeof", "len", "delete"):
		return p.parseBuiltin()
	}

	left := p.parseOperand()
	if left == nil {
		return nil
	}

	// TODO give left to these calls.
	switch curr := p.next(); {
	case curr.Matches("["):
		return p.parseIndex(left)
	case curr.Matches("("):
		return p.parseCall(left)
	}

	return left
}

func (p *parser) parseLeft() *ExpressionNode {
	if expr := p.parsePrimaryExpr(); expr != nil {
		return expr
	}
	return p.parseUnaryExpr()
}

var opPrec = map[string]int{
	"*": 5,
	"/": 5,
	"%": 5,

	"+": 4,
	"-": 4,

	"==": 3,
	"!=": 3,
	"<":  3,
	"<=": 3,
	">":  3,
	">=": 3,

	"&&": 2,

	"||": 1,
}

func getOpPrec(op string) int {
	if prec, ok := opPrec[op]; ok {
		return prec
	}
	return -1
}

func (p *parser) parsePrec(lastPrec int, left *ExpressionNode) *ExpressionNode {
	for p.hasNext() {
		prec := getOpPrec(p.next().Value)
		if prec < lastPrec {
			return left
		}

		// FIXME.
		if !p.hasNext() {
			return left
		}

		// next op is not a binary
		if _, ok := opPrec[p.next().Value]; !ok {
			return left
		}

		op := p.consume()
		right := p.parsePrimaryExpr()
		if right == nil {
			return nil
		}

		if !p.hasNext() {
			return &ExpressionNode{
				Kind:                 BinaryExpression,
				BinaryExpressionNode: &BinaryExpressionNode{left, op.Value, right},
			}
		}

		nextPrec := getOpPrec(p.next().Value)
		if prec < nextPrec {
			right = p.parsePrec(prec+1, right)
			if right == nil {
				return nil
			}
		}

		left = &ExpressionNode{
			Kind: BinaryExpression,
			BinaryExpressionNode: &BinaryExpressionNode{
				left, op.Value, right,
			},
		}
	}

	return left
}

func (p *parser) parseAssign(left *ExpressionNode) *ExpressionNode {
	if !p.hasNext() {
		return nil
	}

	// TODO check if op is valid assign op

	start := p.pos
	op := p.consume()

	right := p.parseExpression()
	if right == nil {
		p.error(api.NewParseError("expression after assignment operator", start, p.pos))
	}

	return &ExpressionNode{
		Kind: AssignStatement,
		AssignStatementNode: &AssignStatementNode{
			left, op.Value, right,
		},
	}
}

func (p *parser) parseDotList(left *ExpressionNode) *ExpressionNode {
	start := p.pos
	list := []*ExpressionNode{}

	list = append(list, left)

	for p.hasNext() && p.next().Matches(".") {
		p.expect(".")
		val := p.parseExpression()
		if val == nil {
			p.error(api.NewParseError("expression in dot-list", start, p.pos))
		}
		list = append(list, val)
	}

	return &ExpressionNode{
		Kind:               PathExpression,
		PathExpressionNode: &PathExpressionNode{list},
	}
}

func (p *parser) parseExpression() *ExpressionNode {
	left := p.parseLeft()
	if left == nil {
		return nil
	}

	if p.next().Matches(".") {
		return p.parseDotList(left)
	}

	// FIXME
	if p.next().Matches("=", "+=", "-=", "*=", "/=") {
		return p.parseAssign(left)
	}

	if _, ok := opPrec[p.next().Value]; ok {
		return p.parsePrec(0, left)
	}
	return left
}

func (p *parser) parseExpressionStatement() *ParseTreeNode {
	expr := p.parseExpression()
	if expr != nil {
		return &ParseTreeNode{
			Kind:                    ExpressionStatement,
			ExpressionStatementNode: expr,
		}
	}

	return nil
}

func (p *parser) parseNode() *ParseTreeNode {
	start := p.pos

	res := &ParseTreeNode{}

	switch curr := p.next(); {
	case curr.Matches(struc):
		res.StructureDeclaration = p.parseStructureDeclaration()
		res.Kind = StructureDeclStatement
	case curr.Matches(trait):
		res.TraitDeclaration = p.parseTraitDeclaration()
		res.Kind = TraitDeclStatement
	case curr.Matches(impl):
		res.ImplDeclaration = p.parseImplDeclaration()
		res.Kind = ImplDeclStatement
	case curr.Matches(fn):
		res.FunctionDeclaration = p.parseFunctionDeclaration()
		res.Kind = FunctionDeclStatement

	case curr.Matches(mut):
		res = p.parseMut()
		p.expect(";")

	case curr.Matches(let):
		res = p.parseLet()
		p.expect(";")

	default:
		p.error(api.NewUnimplementedError(p.next().Value, start, p.pos))
	}

	return res
}

func parseTokenStream(stream []Token) ([]*ParseTreeNode, []api.CompilerError) {
	p := &parser{stream, 0, []api.CompilerError{}}
	nodes := []*ParseTreeNode{}
	for p.hasNext() {
		if node := p.parseNode(); node != nil {
			nodes = append(nodes, node)
		}
	}
	return nodes, p.errors
}
