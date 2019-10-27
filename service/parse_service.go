package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/krug-lang/caasper/entity"
	"github.com/krug-lang/caasper/front"
	"net/http"
)

func Parse(c *gin.Context) {
	var parseReq entity.ParseRequest
	if err := c.BindJSON(&parseReq); err != nil {
		panic(err)
	}

	var stream []front.Token
	if err := jsoniter.Unmarshal([]byte(parseReq.Input), &stream); err != nil {
		panic(err)
	}

	nodes, errors := front.ParseTokenStream(stream)

	jsonNodes, err := jsoniter.MarshalIndent(nodes, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonNodes),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}


func DirectiveParser(c *gin.Context) {
	var directiveReq entity.DirectiveParseRequest
	if err := c.BindJSON(&directiveReq); err != nil {
		panic(err)
	}

	var stream []front.Token
	if err := jsoniter.Unmarshal([]byte(directiveReq.Input), &stream); err != nil {
		panic(err)
	}

	nodes, errors := front.ParseDirectives(stream)

	jsonNodes, err := jsoniter.MarshalIndent(nodes, "", "  ")
	if err != nil {
		panic(err)
	}

	resp := entity.KrugResponse{
		Data:   string(jsonNodes),
		Errors: errors,
	}
	c.JSON(http.StatusOK, &resp)
}
