package api

// frontend

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

// sema

type BuildTypeMapRequest struct{}

type BuildScopeMapRequest struct {
	IRModule string `json:"ir_module"`
}

// borrowcheck

type BorrowCheckRequest struct {
	IRModule string `json:"ir_module"`
	ScopeMap string `json:"scope_map"`
}

// ir build

type IRBuildRequest struct {
	TreeNodes string `json:"tree_nodes"`
}

// codegen

type CodeGenerationRequest struct {
	// IRModule is the string json object for
	// the ir module to generate code for.
	IRModule string `json:"ir_module"`

	// IndentationLevel specifies the amount of indentation
	// for the generated C code, e.g. 4 for a 4 space tab.
	TabSize int `json:"tab_size"`

	// Minify is a flag that specifies whether or not to
	// minify the generated c code
	// this mode simply strips newlines from blocks and
	// structs, etc.
	Minify bool `json:"minify"`
}

type TypeResolveRequest struct{}
type SymbolResolveRequest struct{}

type KrugResponse struct {
	Data   string          `json:"data"`
	Errors []CompilerError `json:"errors"`
}
