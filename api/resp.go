package api

import (
	"encoding/gob"
)

func init() {
	gob.Register(CompilerError{})
	gob.Register(KrugRequest{})
	gob.Register(KrugResponse{})
}

type CompilerError struct {
	Title string
	Desc  string
	Fatal bool
}

type KrugResponse struct {
	Data   []byte
	Errors []CompilerError
}

type KrugRequest struct {
	Data []byte
}
