package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/back"
	"github.com/krug-lang/caasper/entity"
	"github.com/krug-lang/caasper/ir"
	"net/http"
)

func Gen(c *gin.Context) {
	var codeGenReq entity.CodeGenerationRequest
	if err := c.BindJSON(&codeGenReq); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(codeGenReq.IRModule), &irMod); err != nil {
		panic(err)
	}

	// for now we just return the
	// bytes for one big old c file.
	monoFile, errors := back.Codegen(&irMod, codeGenReq.TabSize, codeGenReq.Minify)

	type generatedCode struct {
		Code string `json:"code"`
	}

	genCode := generatedCode{monoFile}
	genCodeResp, err := jsoniter.Marshal(&genCode)
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(genCodeResp),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}

