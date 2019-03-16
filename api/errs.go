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

func NewUnimplementedError(what string, points ...int) CompilerError {
	return CompilerError{
		Title:       fmt.Sprintf("%s unimplemented", what),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewParseError(expected string, points ...int) CompilerError {
	return CompilerError{
		Title:       fmt.Sprintf("Expected %s", expected),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewUnexpectedToken(curr string, expected string, points ...int) CompilerError {
	return CompilerError{
		Title:       fmt.Sprintf("Expected '%s' but found '%s'", expected, curr),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
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
