package api

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(CompilerError{})
}

type CompilerError struct {
	Title       string
	Desc        string
	Fatal       bool
	CodeContext []int
}

func NewUnresolvedSymbol(name string, points ...int) CompilerError {
	return CompilerError{
		Title:       fmt.Sprintf("Unresolved reference to symbol '%s'", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewSymbolError(name string, points ...int) CompilerError {
	return CompilerError{
		Title:       fmt.Sprintf("A symbol with the name '%s' already exists in this scope", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}
