package front

import "io/ioutil"

type KrugCompilationUnit struct {
	Name string
	Code string
}

func ReadCompUnit(loc string) KrugCompilationUnit {
	code, err := ioutil.ReadFile(loc)
	if err != nil {
		panic(err)
	}
	return KrugCompilationUnit{
		Name: "nil",
		Code: string(code),
	}
}
