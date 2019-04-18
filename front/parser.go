package front

import (
	"github.com/hugobrains/caasper/api"
)

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

	err := api.NewUnexpectedToken(val, p.next().Value, start, p.pos)
	p.error(err)
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

func (p *parser) rewind() {
	p.pos--
}

func (p *parser) consume() (tok Token) {
	tok = p.toks[p.pos]
	p.pos++
	return tok
}

func (p *parser) hasNext() bool {
	return p.pos < len(p.toks)
}
