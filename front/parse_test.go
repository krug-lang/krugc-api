package front

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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