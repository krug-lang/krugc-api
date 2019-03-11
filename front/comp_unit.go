package front

import (
	"encoding/gob"
	"io/ioutil"
	"strings"
)

type KrugCompilationUnit struct {
	Name string
	Code string
}

func (k KrugCompilationUnit) GetLine(fst, snd int) string {
	var startLine int
	for startLine = fst; k.Code[startLine] != '\n'; startLine-- {
		// nop!
	}

	var endLine int
	for endLine = snd; k.Code[endLine] != '\n'; endLine++ {
		// nop!
	}

	return strings.TrimSpace(k.Code[startLine+1 : endLine-1])
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
