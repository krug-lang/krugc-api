package front

import (
	"bytes"
	"encoding/gob"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/krugc-api/api"
)

func Parse(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var stream TokenStream
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&stream)

	parseTree, errors := parseTokenStream(&stream)

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	enc.Encode(&parseTree)

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
		Errors: errors,
	}
	c.JSON(200, &resp)
}

func Tokenize(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var sourceFile KrugCompilationUnit
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&sourceFile)

	tokens, errors := tokenizeInput(sourceFile.Code)

	stream := TokenStream{tokens}

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	enc.Encode(&stream)

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
		Errors: errors,
	}
	c.JSON(200, &resp)
}
