package api

import (
	"fmt"
)

type CompilerError struct {
	ErrorCode   int    `json:"error_code"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Fatal       bool   `json:"fatal"`
	CodeContext []int  `json:"code_context"`
}

func NewDirectiveParseError(what string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   1,
		Title:       fmt.Sprintf(what),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewUnimplementedError(phase string, what string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   2,
		Title:       fmt.Sprintf("phase(%s): %s unimplemented", phase, what),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewParseError(expected string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   3,
		Title:       fmt.Sprintf("Expected %s", expected),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewUnexpectedToken(curr string, expected string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   4,
		Title:       fmt.Sprintf("Expected '%s' but found '%s'", expected, curr),
		Desc:        "",
		Fatal:       true,
		CodeContext: points,
	}
}

func NewUnresolvedSymbol(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   5,
		Title:       fmt.Sprintf("Unresolved reference to symbol '%s'", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewSymbolError(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   6,
		Title:       fmt.Sprintf("A symbol with the name '%s' already exists in this scope", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewMovedValueError(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   7,
		Title:       fmt.Sprintf("Use of moved value '%s'", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewMutabilityError(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   8,
		Title:       fmt.Sprintf("Attempting to modify constant value '%s'", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}

func NewUnusedFunction(name string, points ...int) CompilerError {
	return CompilerError{
		ErrorCode:   9,
		Title:       fmt.Sprintf("Function '%s' is not used", name),
		Desc:        "",
		Fatal:       false,
		CodeContext: points,
	}
}
