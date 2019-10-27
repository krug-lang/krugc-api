package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/entity"
	"github.com/krug-lang/caasper/front"
	"io/ioutil"
	"net/http"
	"strings"
)

// Tokenize is the route that handles tokenisation of files.
// It takes api.LexerRequest as the input, containing either
// the code of the file, or the path which must be prefixed
// with an @ symbol.
// If the @ symbol is present, the file is loaded from the
// absolute path provided.
//
// There should be some restrictions on this perhaps... otherwise
// people could 'lex' password files or something.
// Maybe the files must end with a '.krug' extension?
func Tokenize(c *gin.Context) {
	var lexReq entity.LexerRequest
	if err := c.BindJSON(&lexReq); err != nil {
		panic(err)
	}

	code := lexReq.Input
	if len(code) != 0 && code[0] == '@' {
		filePath := strings.Split(code, "@")[1]
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		code = string(data)
	}

	tokens, errors := front.TokenizeInput(code, true)

	jsonResp, err := jsoniter.MarshalIndent(tokens, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonResp),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}

func Comments(c *gin.Context) {
	var commentReq entity.CommentsRequest
	if err := c.BindJSON(&commentReq); err != nil {
		panic(err)
	}

	tokens, errors := front.TokenizeInput(commentReq.Input, false)

	result := []front.Token{}

	// all of the comment tokens.
	for _, tok := range tokens {
		switch tok.Kind {
		case front.SingleLineComment:
			fallthrough
		case front.MultiLineComment:
			result = append(result, tok)
		}
	}

	jsonResp, err := jsoniter.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonResp),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}


