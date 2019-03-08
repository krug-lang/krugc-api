package back

import (
	"fmt"
	"strings"

	"github.com/krug-lang/ir"
)

type emitter struct {
	decl   string
	source string
	target *string
}

func (e *emitter) retarget(to *string) {
	e.target = to
}

func (e *emitter) writeln(f string, d ...interface{}) {
	e.write(f+"\n", d...)
}

func (e *emitter) writetln(num int, f string, d ...interface{}) {
	const TabSize = 4
	tabs := strings.Repeat(" ", TabSize*num)
	e.write(tabs+f+"\n", d...)
}

func (e *emitter) write(f string, d ...interface{}) {
	*e.target += fmt.Sprintf(f, d...)
}

func (e *emitter) writeType(typ ir.Type) string {
	switch t := typ.(type) {

	// compile to uint8_t, uint16_t, etc.
	case *ir.IntegerType:
		var signed rune
		if !t.Signed {
			signed = 'u'
		}
		return fmt.Sprintf("%cint%d_t", signed, t.Width)
	case *ir.FloatingType:
		if t.Width == 32 {
			return "float"
		}
		return "double"
	}

	panic("unhandled type")
}

func (e *emitter) buildAlloca(a *ir.Alloca) {
	aType := e.writeType(a.Type)
	e.writeln("%s %s = malloc(sizeof(*%s));", aType, a.Name, a.Name)
}

func (e *emitter) buildExpr(l ir.Value) string {
	switch val := l.(type) {
	case *ir.IntegerValue:
		return val.RawValue.String()

	case *ir.BinaryExpression:
		lh := e.buildExpr(val.LHand)
		rh := e.buildExpr(val.RHand)
		return fmt.Sprintf("(%s%s%s)", lh, val.Op, rh)
	}

	panic("oh no")
}

func (e *emitter) buildLocal(l *ir.Local) {
	aType := e.writeType(l.Type)
	modifier := ""
	if !l.Mutable {
		modifier = "const "
	}

	valueCode := ";"
	if l.Val != nil {
		valueCode = fmt.Sprintf(" = %s;", e.buildExpr(l.Val))
	}

	e.writetln(1, "%s%s %s%s", modifier, aType, l.Name, valueCode)
}

func (e *emitter) buildRet(r *ir.Return) {
	res := ";"
	if r.Val != nil {
		res = fmt.Sprintf(" %s;", e.buildExpr(r.Val))
	}
	e.writetln(1, "return%s", res)
}

func (e *emitter) buildInstr(i ir.Instruction) {
	switch instr := i.(type) {
	case *ir.Alloca:
		e.buildAlloca(instr)
		return
	case *ir.Local:
		e.buildLocal(instr)
		return
	case *ir.Block:
		e.buildBlock(instr)
		return
	case *ir.Return:
		e.buildRet(instr)
		return
	}

	panic("unhandled instr")
}

func (e *emitter) buildBlock(b *ir.Block) {
	e.writeln("{")
	for _, instr := range b.Instr {
		e.buildInstr(instr)
	}
	e.writeln("}")
}

func (e *emitter) emitFunc(fn *ir.Function) {
	writeArgList := func(fn *ir.Function) string {
		var argList string
		idx := 0
		for name, t := range fn.Param {
			if idx != 0 {
				argList += ", "
			}
			argList += fmt.Sprintf("%s %s", e.writeType(t), name)
			idx++
		}
		return argList
	}

	typ := e.writeType(fn.ReturnType)

	// write prototype to the decl part
	e.retarget(&e.decl)
	argList := writeArgList(fn)
	e.writeln("%s %s(%s);", typ, fn.Name, argList)

	// write definition to the source part.
	e.retarget(&e.source)
	e.writeln("%s %s(%s)", typ, fn.Name, argList)
	e.buildBlock(fn.Body)
}

func codegen(mod *ir.Module) []byte {
	e := &emitter{}
	e.retarget(&e.source)

	for _, fn := range mod.Functions {
		e.emitFunc(fn)
	}

	return []byte(e.decl + e.source)
}
