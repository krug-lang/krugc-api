package api

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(CompilerError{})
}

type CompilerError struct {
	Title string
	Desc  string
	Fatal bool
}

func NewSymbolError(name string) CompilerError {
	return CompilerError{
		Title: fmt.Sprintf("A symbol with the name '%s' already exists in this scope", name),
		Desc:  "",
		Fatal: false,
	}
}
