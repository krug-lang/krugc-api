package entity

import "github.com/krug-lang/caasper/api"

type KrugResponse struct {
	Data   string          `json:"data"`
	Errors []api.CompilerError `json:"errors"`
}

