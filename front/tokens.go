package front

type TokenType string

const (
	Identifier TokenType = "iden"
)

type Token struct {
	Value string
	Kind  TokenType
}

type TokenStream struct {
	Tokens []*Token
}
