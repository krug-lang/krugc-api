package front

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
)

func Parse(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var stream []Token
	if err := jsoniter.Unmarshal([]byte(krugReq.Data), &stream); err != nil {
		panic(err)
	}

	nodes, errors := parseTokenStream(stream)

	jsonNodes, err := jsoniter.MarshalIndent(nodes, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonNodes),
		Errors: errors,
	}
	c.JSON(200, &resp)
}

func Tokenize(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	code := string(krugReq.Data)

	tokens, errors := tokenizeInput(code)

	jsonResp, err := jsoniter.MarshalIndent(tokens, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonResp),
		Errors: errors,
	}
	c.JSON(200, &resp)
}
