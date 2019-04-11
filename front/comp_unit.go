package front

import (
	"io/ioutil"
	"strings"
	"unicode"
)

type KrugCompilationUnit struct {
	Name string
	Code string
}

// GetLine will return the line of the given position, as well
// as the position underlined as a second (optional) result.
func (k KrugCompilationUnit) GetLine(fst, snd int) (string, string) {
	var startLine int
	for startLine = fst; k.Code[startLine] != '\n'; startLine-- {
		// nop!
	}

	var endLine int
	for endLine = snd; k.Code[endLine] != '\n'; endLine++ {
		// nop!
	}

	// extract the result and make a note of its left.
	result := []rune(k.Code[startLine+1 : endLine-1])
	oldLength := len(result)

	// remove all cruft from the LEFT of the string
	result = []rune(strings.TrimLeftFunc(string(result), unicode.IsSpace))

	// work out how many characters we've cut off
	offset := oldLength - len(result)

	// work out the length of the token to underline
	length := snd - fst

	// calculate where the underline should
	// start, take into account the offset from
	// removing the spaces at the start.
	start := fst - (startLine + 1) - offset

	// write the underline.
	underlined := []rune(strings.Repeat(" ", len(result)))
	for i := 0; i < length; i++ {
		underlined[start+i] = '^'
	}

	// run trimspace on the final result to remove any
	// rubbish from the right of the string
	return strings.TrimSpace(string(result)), string(underlined)
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
