package front

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/krug-serv/api"
)

func Parse(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var stream TokenStream
	/*
		var stream TokenStream
		pCache := bytes.NewBuffer(krugReq.Data)
		decCache := gob.NewDecoder(pCache)
		decCache.Decode(&stream)
	*/

	parseTree, errors := parseTokenStream(&stream)

	jsonParseTree, err := jsoniter.MarshalIndent(parseTree, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonParseTree),
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
	sourceFile.Code = string(krugReq.Data)

	/*
		pCache := bytes.NewBuffer(krugReq.Data)
		decCache := gob.NewDecoder(pCache)
		decCache.Decode(&sourceFile)
	*/

	tokens, errors := tokenizeInput(sourceFile.Code)

	stream := TokenStream{tokens}

	jsonResp, err := jsoniter.MarshalIndent(stream, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonResp),
		Errors: errors,
	}
	c.JSON(200, &resp)
}
