package back

import (
	"bytes"
	"encoding/gob"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/ir"
	"github.com/krug-lang/krugc-api/api"
)

func Gen(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod *ir.Module
	pCache := bytes.NewBuffer(krugReq.Data)
	decCache := gob.NewDecoder(pCache)
	decCache.Decode(&irMod)

	// for now we just return the
	// bytes for one big old c file.
	monoFile := codegen(irMod)

	resp := api.KrugResponse{
		Data: monoFile,
	}
	c.JSON(200, &resp)
}
