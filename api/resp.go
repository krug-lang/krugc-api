package api

import (
	"encoding/gob"
)

func init() {
	gob.Register(KrugRequest{})
	gob.Register(KrugResponse{})
}

type KrugResponse struct {
	Data   []byte
	Errors []CompilerError
}

type KrugRequest struct {
	Data []byte
}
