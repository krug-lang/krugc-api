package front

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func tok(value string, kind TokenType) Token {
	return Token{value, kind, []int{}}
}

func tokensMatch(t *testing.T, a, b Token) {
	assert.Equal(t, a.Kind, b.Kind)
	assert.Equal(t, a.Value, b.Value)
}

func tokenSetMatches(t *testing.T, a []Token, b []Token) {
	assert.Equal(t, len(b), len(a))
	for idx, tok := range a {
		assert.Equal(t, tok.Value, b[idx].Value)
		assert.Equal(t, tok.Kind, b[idx].Kind)
	}
}

func TestCommentIgnoreToggle(t *testing.T) {
	t.Log("Checking if comments are ignored")
	{
		tokens, errs := tokenizeInput("// this is a comment", true)
		assert.Empty(t, tokens)
		assert.Empty(t, errs)
	}

	t.Log("Checking if comments are not ignored")
	{
		tokens, errs := tokenizeInput("// this is a comment", false)
		assert.NotEmpty(t, tokens)
		assert.Empty(t, errs)
	}
}

func TestLexEmptyFunction(t *testing.T) {
	tokens, errs := tokenizeInput("fn main() int { }", false)
	assert.Empty(t, errs)
	tokenSetMatches(t, tokens, []Token{
		tok("fn", Identifier),
		tok("main", Identifier),
		tok("(", Symbol),
		tok(")", Symbol),
		tok("int", Identifier),
		tok("{", Symbol),
		tok("}", Symbol),
	})
}

func TestLoopingConstructs(t *testing.T) {
	t.Log("Testing infinite loop")
	tokens, errs := tokenizeInput("loop { }", false)
	assert.Empty(t, errs)
	tokenSetMatches(t, tokens, []Token{
		tok("loop", Identifier),
		tok("{", Symbol),
		tok("}", Symbol),
	})

	t.Log("Testing while loop")
	tokens, errs = tokenizeInput(`while i < 100; i = i + 1 { printf("%d\n") }`, false)
	assert.Empty(t, errs)
	tokenSetMatches(t, tokens, []Token{
		tok("while", Identifier),
		tok("i", Identifier),
		tok("<", Symbol),
		tok("100", Number),
		tok(";", Symbol),
		tok("i", Identifier),
		tok("=", Symbol),
		tok("i", Identifier),
		tok("+", Symbol),
		tok("1", Number),
		tok("{", Symbol),
		tok("printf", Identifier),
		tok("(", Symbol),
		tok(`"%d\n"`, String),
		tok(")", Symbol),
		tok("}", Symbol),
	})
}

func TestLexVariable(t *testing.T) {
	tokens, errs := tokenizeInput("let x int = 3;", false)
	assert.Empty(t, errs)

	tokenSetMatches(t, tokens, []Token{
		tok("let", Identifier),
		tok("x", Identifier),
		tok("int", Identifier),
		tok("=", Symbol),
		tok("3", Number),
		tok(";", Symbol),
	})
}
