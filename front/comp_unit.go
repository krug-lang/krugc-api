package front

import (
	"encoding/gob"
	"io/ioutil"
)

type KrugCompilationUnit struct {
	Name string
	Code string
}

func init() {
	gob.Register(KrugCompilationUnit{})
}

func ReadCompUnit(loc string) KrugCompilationUnit {
	code, err := ioutil.ReadFile(loc)
	if err != nil {
		panic(err)
	}
	return KrugCompilationUnit{
		Name: "foopa",
		Code: string(code),
	}
}
