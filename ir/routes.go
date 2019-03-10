package ir

import (
	"bytes"
	"encoding/gob"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/krugc-api/api"
	"github.com/krug-lang/krugc-api/front"
)

// Build takes the given []ParseTree's
// and builds a SINGLE ir module from them.
func Build(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var trees []front.ParseTree
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&trees)

	irModule, errors := build(trees)

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	enc.Encode(&irModule)

	resp := api.KrugResponse{
		Data:   buff.Bytes(),
		Errors: errors,
	}
	c.JSON(200, &resp)
}
