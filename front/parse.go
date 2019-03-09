package front

import (
	"fmt"
	"strings"
)

type parser struct {
	toks []Token
	pos  int
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
	if p.hasNext() {
		if tok = p.consume(); tok.Matches(val) {
			return tok
		}
	}
	panic(fmt.Sprintf("Expected '%s', got '%s'", val, tok))
}

func (p *parser) expectKind(kind TokenType) (tok Token) {
	if tok = p.consume(); tok.Kind == kind {
		return tok
	}
	panic(fmt.Sprintf("Expected '%s', got '%s'", kind, tok))
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
	p.expect("*")
	base := p.parseType()
	if base == nil {
		panic("no type after ptr")
	}
	return NewPointerType(base)
}

func (p *parser) parseArrayType() *ArrayType {
	p.expect("[")
	base := p.parseType()
	if base == nil {
		panic("no type after array")
	}

	// TODO length expect(";"), parseConstant.

	p.expect("]")
	return NewArrayType(base)
}

func (p *parser) parseUnresolvedType() *UnresolvedType {
	name := p.expectKind(Identifier)
	return NewUnresolvedType(name.Value)
}

func (p *parser) parseType() TypeNode {
	switch curr := p.next(); {
	case curr.Matches("*"):
		return p.parsePointerType()
	case curr.Matches("["):
		return p.parseArrayType()
	case curr.Kind == Identifier:
		return p.parseUnresolvedType()
	default:
		panic(fmt.Sprintf("unimplemented type %s", curr))
	}
}

func (p *parser) parseStructureDeclaration() *StructureDeclaration {
	p.expect("struct")
	name := p.consume().Value

	fields := map[string]TypeNode{}

	p.expect("{")
	for p.hasNext() {
		if p.next().Matches("}") {
			break
		}

		name := p.expectKind(Identifier)

		typ := p.parseType()
		if typ == nil {
			panic("no type!")
		}
		fields[name.Value] = typ

		// trailing commas are enforced.
		p.expect(",")
	}
	p.expect("}")

	return NewStructureDeclaration(name, fields)
}

func (p *parser) parseFunctionPrototypeDeclaration() *FunctionPrototypeDeclaration {
	p.expect("func")
	name := p.expectKind(Identifier).Value

	args := map[string]TypeNode{}

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
			panic("expected type in func proto")
		}

		args[name.Value] = typ
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
	p.expect("mut")
	name := p.expectKind(Identifier)

	var typ TypeNode
	if !p.next().Matches("=") {
		typ = p.parseType()
		if typ == nil {
			panic("expected type or assignment")
		}
	}

	var val ExpressionNode
	if p.next().Matches("=") {
		p.expect("=")

		val = p.parseExpression()
		if val == nil {
			panic("expected assignment")
		}
	}

	if val == nil && typ == nil {
		panic("expected value or type in mutable statement")
	}

	return NewMutableStatement(name.Value, typ, val)
}

// let is a constant variable.
func (p *parser) parseLet() StatementNode {
	p.expect("let")
	name := p.expectKind(Identifier)

	var typ TypeNode
	if !p.next().Matches("=") {
		typ = p.parseType()
		if typ == nil {
			panic("expected type or assignment")
		}
	}

	// let MUST have a value assigned.
	p.expect("=")
	value := p.parseExpression()
	return NewLetStatement(name.Value, typ, value)
}

func (p *parser) parseReturn() StatementNode {
	if !p.next().Matches("return") {
		return nil
	}
	p.expect("return")

	var res ExpressionNode
	if !p.next().Matches(";") {
		res = p.parseExpression()
		if res == nil {
			panic("return expected expression or semi-colon")
		}
	}

	return NewReturnStatement(res)
}

func (p *parser) parseSemicolonStatement() StatementNode {
	switch curr := p.next(); {
	case curr.Matches("mut"):
		return p.parseMut()
	case curr.Matches("let"):
		return p.parseLet()
	case curr.Matches("return"):
		return p.parseReturn()
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

	// the first if.
	p.expect("if")
	expr := p.parseExpression()
	if expr == nil {
		panic("wanted condition")
	}

	block := p.parseStatBlock()
	if block == nil {
		panic("expected block after condition")
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
				panic("expected cond in else if")
			}

			body := p.parseStatBlock()
			if body == nil {
				panic("else if expected block")
			}
			elses = append(elses, NewElseIfStatement(cond, body))
		} else {
			// TODO we could easily check if else has been set before
			// or should we do this in a later pass during sema analyis?

			p.expect("else")
			elseBlock = p.parseStatBlock()
			if elseBlock == nil {
				panic("else expected block")
			}
		}
	}

	return NewIfStatement(expr, block, elses, elseBlock)
}

func (p *parser) parseWhileLoop() StatementNode {
	if !p.next().Matches("while") {
		return nil
	}
	p.expect("while")
	val := p.parseExpression()
	if val == nil {
		panic("while expects condition")
	}

	var post ExpressionNode
	if p.next().Matches(";") {
		p.expect(";")
		post = p.parseExpression()
		if post == nil {
			panic("expected expression after ; in while loop")
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
	name := p.expectKind(Identifier).Value

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

	return NewImplDeclaration(name)
}

func (p *parser) parseTraitDeclaration() *TraitDeclaration {
	p.expect("trait")

	name := p.expectKind(Identifier).Value

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

	op := p.consume()
	right := p.parseLeft()
	if right == nil {
		panic("error, unary expr is null!")
	}

	return NewUnaryExpression(op.Value, right)
}

func (p *parser) parseOperand() ExpressionNode {
	if !p.hasNext() {
		return nil
	}

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
		return NewVariableReference(curr.Value)

	case EndOfFile:
		return nil

	default:
		panic(fmt.Sprintf("unhandled thingy %s", curr))
	}
}

func (p *parser) parseBuiltin() ExpressionNode {
	builtin := p.expectKind(Identifier)
	p.expect("!")
	typ := p.parseType()
	if typ == nil {
		panic("Expected type after builtin")
	}
	return NewBuiltinExpression(builtin.Value, typ)
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
		// parseIndex
		return nil
	case curr.Matches("("):
		// parseCall
		return nil
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

	op := p.consume()

	right := p.parseExpression()
	if right == nil {
		panic("assign expected expression after assignment operator")
	}

	return NewAssignmentStatement(left, op.Value, right)
}

func (p *parser) parseExpression() ExpressionNode {
	left := p.parseLeft()
	if left == nil {
		return nil
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

	panic(fmt.Sprintf("unhandled %s", p.next()))
}

func parseTokenStream(stream *TokenStream) ParseTree {
	p := &parser{stream.Tokens, 0}
	nodes := []ParseTreeNode{}
	for p.hasNext() {
		if node := p.parseNode(); node != nil {
			nodes = append(nodes, node)
		}
	}
	return ParseTree{nodes}
}
