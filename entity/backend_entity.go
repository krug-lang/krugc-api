package entity

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
