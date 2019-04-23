package api

import (
	"fmt"
)

// TODO ensure that error codes arent
// entered manually, or at least, they are
// consistent.

type CompilerError struct {
	ErrorCode   int    `json:"error_code"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Fatal       bool   `json:"fatal"`
	CodeContext []int  `json:"code_context"`
}

func NewDirectiveParseError(what string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0001,
		Title:       fmt.Sprintf(what),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewUnimplementedError(what string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0002,
		Title:       fmt.Sprintf("%s unimplemented", what),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewParseError(expected string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0003,
		Title:       fmt.Sprintf("Expected %s", expected),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewUnexpectedToken(curr string, expected string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0004,
		Title:       fmt.Sprintf("Expected '%s' but found '%s'", expected, curr),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewUnresolvedSymbol(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0005,
		Title:       fmt.Sprintf("Unresolved reference to symbol '%s'", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewSymbolError(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0006,
		Title:       fmt.Sprintf("A symbol with the name '%s' already exists in this scope", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewMovedValueError(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0007,
		Title:       fmt.Sprintf("Use of moved value '%s'", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewMutabilityError(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   0007,
		Title:       fmt.Sprintf("Attempting to modify constant value '%s'", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}
