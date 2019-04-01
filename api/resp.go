package api

import (
	"encoding/gob"
)

func init() {
	gob.Register(KrugRequest{})
	gob.Register(KrugResponse{})
}

type KrugResponse struct {
	Data   string          `json:"data"`
	Errors []CompilerError `json:"errors"`
}

type KrugRequest struct {
	Data string `json:"data"`
}
