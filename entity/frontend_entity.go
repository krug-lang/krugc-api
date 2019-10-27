package entity

type LexerRequest struct {
	Input string `json:"input"`
}

type CommentsRequest struct {
	Input string `json:"input"`
}

type DirectiveParseRequest struct {
	Input string `json:"input"`
}

type ParseRequest struct {
	Input string `json:"input"`
}
