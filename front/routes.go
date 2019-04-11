package front

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
)

func Comments(c *gin.Context) {
	var commentReq api.CommentsRequest
	if err := c.BindJSON(&commentReq); err != nil {
		panic(err)
	}

	tokens, errors := tokenizeInput(commentReq.Input, false)

	result := []Token{}

	// all of the comment tokens.
	for _, tok := range tokens {
		switch tok.Kind {
		case SingleLineComment:
			fallthrough
		case MultiLineComment:
			result = append(result, tok)
		}
	}

	jsonResp, err := jsoniter.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := api.KrugResponse{
		Data:   string(jsonResp),
		Errors: errors,
	}
	c.JSON(200, &resp)
}

func Parse(c *gin.Context) {
	var parseReq api.ParseRequest
	if err := c.BindJSON(&parseReq); err != nil {
		panic(err)
	}

	var stream []Token
	if err := jsoniter.Unmarshal([]byte(parseReq.Input), &stream); err != nil {
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
	var lexReq api.LexerRequest
	if err := c.BindJSON(&lexReq); err != nil {
		panic(err)
	}

	tokens, errors := tokenizeInput(lexReq.Input, true)

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
