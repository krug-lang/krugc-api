package front

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hugobrains/caasper/api"
)

type directiveParser struct {
	parser
}

func (p *directiveParser) hasNext() bool {
	return p.pos < len(p.toks) && len(p.errors) == 0
}

type parseFn func(p *directiveParser) *Directive

func (p *directiveParser) parseArgumentList() []value {
	p.expect("(")

	vals := []value{}

	idx := 0
	for p.hasNext() && !p.next().Matches(")") {
		if idx != 0 {
			p.expect(",")
		}

		var kind valueKind
		tok := p.consume()
		switch tok.Kind {
		case String:
			kind = stringValue
		case Number:
			if strings.Index(tok.Value, ".") == -1 {
				kind = integerValue
			} else {
				kind = floatingValue
			}
		case Char:
			kind = characterValue
		default:
			panic(fmt.Sprintf("unhandled directive value %s", tok.Value))
		}

		vals = append(vals, value{kind, tok})
		idx++
	}

	p.expect(")")

	return vals
}

func parseInclude(p *directiveParser) *Directive {
	start := p.pos

	args := p.parseArgumentList()
	if args == nil {
		p.error(api.NewDirectiveParseError("expected argument list", start, p.pos))
		return nil
	}

	if len(args) != 1 {
		// TODO error: not enough arguments supplied.
		p.error(api.NewDirectiveParseError("not enough arguments supplied", start, p.pos))
		return nil
	}

	// TODO improve the type checking for directive args.
	if args[0].kind != stringValue {
		p.error(api.NewDirectiveParseError("include should have on parameter of type 'string'", start, p.pos))
		return nil
	}

	path := args[0].value.Value
	return &Directive{
		Kind:             Include,
		IncludeDirective: &IncludeDirective{path},
	}
}

func parseLink(p *directiveParser) *Directive {
	args := p.parseArgumentList()
	if args == nil {
		// TODO(ERROR)
		return nil
	}

	flags := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		flags[i] = args[i].value.Value
	}

	return &Directive{
		Kind:          Link,
		LinkDirective: &LinkDirective{flags},
	}
}

func parseNoMangle(p *directiveParser) *Directive {
	if !p.next().Matches("}") {
		return nil
	}

	return &Directive{
		Kind:              NoMangle,
		NoMangleDirective: &NoMangleDirective{},
	}
}

func parseAlign(p *directiveParser) *Directive {
	args := p.parseArgumentList()
	if args == nil {
		// TODO(ERROR)
		return nil
	}

	if len(args) != 1 {
		// TODO(ERROR)
		return nil
	}

	val := args[0].value.Value
	alignment, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		// TODO(ERROR)
		panic(err)
	}

	return &Directive{
		Kind:           Align,
		AlignDirective: &AlignDirective{alignment},
	}
}

func parsePacked(p *directiveParser) *Directive {
	if !p.next().Matches("}") {
		return nil
	}

	return &Directive{
		Kind:            Packed,
		PackedDirective: &PackedDirective{},
	}
}

func parseClang(p *directiveParser) *Directive {
	if !p.next().Matches("}") {
		return nil
	}

	return &Directive{
		Kind:           Clang,
		ClangDirective: &ClangDirective{},
	}
}

func (p *directiveParser) parseDirective() []*Directive {
	p.expect("#")
	p.expect("{")

	luTable := map[string]parseFn{
		"include":   parseInclude,
		"link":      parseLink,
		"no_mangle": parseNoMangle,
		"align":     parseAlign,
		"packed":    parsePacked,
		"clang":     parseClang,
	}

	dirs := []*Directive{}

	idx := 0
	for p.hasNext() {
		if idx != 0 {
			p.expect(",")
		}

		word := p.expectKind(Identifier)
		fmt.Println(word)
		if res, ok := luTable[word.Value]; ok {
			if dir := res(p); dir != nil {
				dirs = append(dirs, dir)
			}
		} else {
			// TODO(ERROR): unrecognize directive
			panic(fmt.Sprintf("no such directive! %s", word))
		}

		if p.next().Matches("}") {
			break
		}

		idx++
	}

	p.expect("}")

	fmt.Println(dirs)

	return dirs
}

func parseDirectives(toks []Token) ([]*Directive, []api.CompilerError) {
	p := &directiveParser{parser{toks, 0, []api.CompilerError{}}}

	nodes := []*Directive{}
	for p.hasNext() {
		if curr := p.next(); curr.Matches("#") {
			nodes = append(nodes, p.parseDirective()...)
		} else {
			p.consume()
		}
	}
	return nodes, p.errors
}
