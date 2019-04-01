package front

import (
	"strings"

	"github.com/krug-lang/server/api"
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

func (p *parser) parsePointerType() *PointerType {
	start := p.pos
	p.expect("^")
	base := p.parseType()
	if base == nil {
		p.error(api.NewParseError("type after pointer", start, p.pos))
		return nil
	}
	return NewPointerType(base)
}

func (p *parser) parseArrayType() *ArrayType {
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
	return NewArrayType(base, size)
}

func (p *parser) parseUnresolvedType() *UnresolvedType {
	name := p.expectKind(Identifier)
	return NewUnresolvedType(name.Value)
}

func (p *parser) parseType() TypeNode {
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
		fields = append(fields, &NamedType{name, typ})

		// trailing commas are enforced.
		p.expect(",")
	}
	p.expect("}")

	return NewStructureDeclaration(name, fields)
}

func (p *parser) parseFunctionPrototypeDeclaration() *FunctionPrototypeDeclaration {
	start := p.pos

	p.expect("func")
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

	fpt := NewFunctionPrototypeDeclaration(name, args)

	if typ := p.parseType(); typ != nil {
		fpt.SetReturnType(typ)
	}

	return fpt
}

// mut x [ type ] [ = val ]
func (p *parser) parseMut() StatementNode {
	start := p.pos

	p.expect("mut")
	name := p.expectKind(Identifier)

	var typ TypeNode
	if !p.next().Matches("=") {
		typ = p.parseType()
		if typ == nil {
			p.error(api.NewParseError("type after assignment", start, p.pos))
		}
	}

	var val ExpressionNode
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

	return NewMutableStatement(name, typ, val)
}

// let is a constant variable.
func (p *parser) parseLet() StatementNode {
	start := p.pos

	p.expect("let")
	name := p.expectKind(Identifier)

	var typ TypeNode
	if !p.next().Matches("=") {
		typ = p.parseType()
		if typ == nil {
			p.error(api.NewParseError("type or assignment", start, p.pos))
		}
	}

	var value ExpressionNode
	if p.next().Matches("=") {
		p.expect("=")
		value = p.parseExpression()
		if value == nil {
			p.error(api.NewParseError("expression in let statement", start, p.pos))
		}
	}
	return NewLetStatement(name, typ, value)
}

func (p *parser) parseReturn() StatementNode {
	if !p.next().Matches("return") {
		return nil
	}
	start := p.pos
	p.expect("return")

	var res ExpressionNode
	if !p.next().Matches(";") {
		res = p.parseExpression()
		if res == nil {
			p.error(api.NewParseError("semi-colon or expression", start, p.pos))
		}
	}

	return NewReturnStatement(res)
}

func (p *parser) parseNext() StatementNode {
	p.expect("next")
	return NewNextStatement()
}

func (p *parser) parseBreak() StatementNode {
	p.expect("break")
	return NewBreakStatement()
}

func (p *parser) parseSemicolonStatement() StatementNode {
	switch curr := p.next(); {
	case curr.Matches("mut"):
		return p.parseMut()
	case curr.Matches("let"):
		return p.parseLet()
	case curr.Matches("return"):
		return p.parseReturn()
	case curr.Matches("next"):
		return p.parseNext()
	case curr.Matches("break"):
		return p.parseBreak()
	}

	return p.parseExpressionStatement()
}

func (p *parser) parseStatBlock() []StatementNode {
	if !p.next().Matches("{") {
		return nil
	}

	stats := []StatementNode{}
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
	return stats
}

func (p *parser) parseIfElseChain() StatementNode {
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

	var elseBlock []StatementNode
	elses := []*ElseIfStatement{}

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
			elses = append(elses, NewElseIfStatement(cond, body))
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

	return NewIfStatement(expr, block, elses, elseBlock)
}

func (p *parser) parseWhileLoop() StatementNode {
	if !p.next().Matches("while") {
		return nil
	}
	start := p.pos

	p.expect("while")
	val := p.parseExpression()
	if val == nil {
		p.error(api.NewParseError("condition after while", start, p.pos))
	}

	var post ExpressionNode
	if p.next().Matches(";") {
		p.expect(";")
		post = p.parseExpression()
		if post == nil {
			p.error(api.NewParseError("step expression in while loop", start, p.pos))
		}
	}

	if block := p.parseStatBlock(); block != nil {
		return NewWhileLoopStatement(val, post, block)
	}
	return nil
}

func (p *parser) parseLoop() StatementNode {
	if !p.next().Matches("loop") {
		return nil
	}
	p.expect("loop")
	if block := p.parseStatBlock(); block != nil {
		return NewLoopStatement(block)
	}
	return nil
}

func (p *parser) parseStatement() StatementNode {
	switch curr := p.next(); {
	case curr.Matches("if"):
		return p.parseIfElseChain()
	case curr.Matches("loop"):
		return p.parseLoop()
	case curr.Matches("while"):
		return p.parseWhileLoop()
	case curr.Matches("{"):
		return NewBlockNode(p.parseStatBlock())
	}

	stat := p.parseSemicolonStatement()
	if stat != nil {
		p.expect(";")
	}
	return stat
}

func (p *parser) parseFunctionDeclaration() *FunctionDeclaration {
	proto := p.parseFunctionPrototypeDeclaration()

	stats := []StatementNode{}

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
	return NewFunctionDeclaration(proto, stats)
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

	return NewImplDeclaration(name, functions)
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

	return NewTraitDeclaration(name, members)
}

func (p *parser) parseUnaryExpr() ExpressionNode {
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

	return NewUnaryExpression(op.Value, right)
}

func (p *parser) parseOperand() ExpressionNode {
	if !p.hasNext() {
		return nil
	}

	start := p.pos
	curr := p.next()

	if curr.Matches("(") {
		p.expect("(")
		expr := p.parseExpression()
		p.expect(")")
		return NewGrouping(expr)
	}

	switch curr := p.consume(); curr.Kind {
	case Number:
		// no dot means it's a whole number.
		if strings.Index(curr.Value, ".") == -1 {
			return NewIntegerConstant(curr.Value)
		}
		return NewFloatingConstant(curr.Value)
	case Identifier:
		return NewVariableReference(curr)

	case String:
		return NewStringConstant(curr.Value)

	case EndOfFile:
		return nil

	default:
		p.error(api.NewUnimplementedError(curr.Value, start, p.pos))
		return nil
	}
}

func (p *parser) parseBuiltin() ExpressionNode {
	start := p.pos
	builtin := p.expectKind(Identifier)
	p.expect("!")
	typ := p.parseType()
	if typ == nil {
		p.error(api.NewParseError("type in builtin", start, p.pos))
	}
	return NewBuiltinExpression(builtin.Value, typ)
}

func (p *parser) parseCall(left ExpressionNode) ExpressionNode {
	start := p.pos

	var params []ExpressionNode

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

	return NewCallExpression(left, params)
}

func (p *parser) parseIndex(left ExpressionNode) ExpressionNode {
	start := p.pos
	p.expect("[")
	val := p.parseExpression()
	if val == nil {
		p.error(api.NewParseError("expression in array index", start, p.pos))
	}
	p.expect("]")
	return NewIndexExpression(left, val)
}

func (p *parser) parsePrimaryExpr() ExpressionNode {
	if !p.hasNext() {
		return nil
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

func (p *parser) parseLeft() ExpressionNode {
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

func (p *parser) parsePrec(lastPrec int, left ExpressionNode) ExpressionNode {
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
			return NewBinary(left, op.Value, right)
		}

		nextPrec := getOpPrec(p.next().Value)
		if prec < nextPrec {
			right = p.parsePrec(prec+1, right)
			if right == nil {
				return nil
			}
		}

		left = NewBinary(left, op.Value, right)
	}

	return left
}

func (p *parser) parseAssign(left ExpressionNode) ExpressionNode {
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

	return NewAssignmentStatement(left, op.Value, right)
}

func (p *parser) parseDotList(left ExpressionNode) ExpressionNode {
	start := p.pos
	list := []ExpressionNode{}

	list = append(list, left)

	for p.hasNext() && p.next().Matches(".") {
		p.expect(".")
		val := p.parseExpression()
		if val == nil {
			p.error(api.NewParseError("expression in dot-list", start, p.pos))
		}
		list = append(list, val)
	}

	return NewPathExpression(list)
}

func (p *parser) parseExpression() ExpressionNode {
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

func (p *parser) parseExpressionStatement() ExpressionNode {
	expr := p.parseExpression()
	if expr != nil {
		return expr
	}

	return nil
}

func (p *parser) parseNode() ParseTreeNode {
	start := p.pos

	switch curr := p.next(); {
	case curr.Matches("struct"):
		return p.parseStructureDeclaration()
	case curr.Matches("trait"):
		return p.parseTraitDeclaration()
	case curr.Matches("impl"):
		return p.parseImplDeclaration()
	case curr.Matches("func"):
		return p.parseFunctionDeclaration()
	}

	p.error(api.NewUnimplementedError(p.next().Value, start, p.pos))
	return nil
}

func parseTokenStream(stream *TokenStream) (ParseTree, []api.CompilerError) {
	p := &parser{stream.Tokens, 0, []api.CompilerError{}}
	nodes := []ParseTreeNode{}
	for p.hasNext() {
		if node := p.parseNode(); node != nil {
			nodes = append(nodes, node)
		}
	}
	return ParseTree{nodes}, p.errors
}
