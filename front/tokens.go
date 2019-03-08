package front

import (
	"encoding/gob"
	"fmt"
	"strings"
)

type TokenType string

const (
	Identifier TokenType = "iden"
	Symbol               = "sym"
	String               = "string"
	Char                 = "char"
	Number               = "num"
	EndOfFile            = "<eof>"
)

type Token struct {
	Value string
	Kind  TokenType
	Span  []int
}

func init() {
	gob.Register(Token{})
	gob.Register(TokenStream{})
}

// Matches returns if the tokens LEXEME
// is identical to the given string
func (t Token) Matches(values ...string) bool {
	for _, value := range values {
		if strings.Compare(t.Value, value) == 0 {
			return true
		}
	}
	return false
}

// Exactly returns if the tokens value and type
// are the exact same as the given value and type.
func (t Token) Exactly(value string, typ TokenType) bool {
	return t.Kind == typ && t.Matches(value)
}

func (t Token) String() string {
	return fmt.Sprintf("[%s](%s)", t.Value, string(t.Kind))
}

func NewToken(lexeme string, kind TokenType, start, end int) Token {
	return Token{
		lexeme, kind, []int{start, end},
	}
}

type TokenStream struct {
	Tokens []Token
}
