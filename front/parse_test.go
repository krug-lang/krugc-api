package front

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMutStatementParses(t *testing.T) {
	input := []Token {
		tok("mut", Identifier),
		tok("x", Identifier),
		tok("=", Symbol),
		tok("5", Number),
		tok(";", Symbol),
	}
	nodes, errs := ParseTokenStream(input)

	mutStat := nodes[0].MutableStatementNode
	assert.NotNil(t, mutStat)
	assert.Equal(t, "x", mutStat.Name.Value)

	assert.NotEmpty(t, nodes)
	assert.Empty(t, errs)
}

func TestLetStatementParses(t *testing.T) {
	input := []Token {
		tok("let", Identifier),
		tok("x", Identifier),
		tok("=", Symbol),
		tok("5", Number),
		tok(";", Symbol),
	}
	nodes, errs := ParseTokenStream(input)

	letStat := nodes[0].LetStatementNode
	assert.NotNil(t, letStat)
	assert.Equal(t, "x", letStat.Name.Value)

	assert.NotEmpty(t, nodes)
	assert.Empty(t, errs)
}