package back

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hugobrains/krug-serv/api"
	"github.com/hugobrains/krug-serv/ir"
)

type emitter struct {
	decl        string
	source      string
	target      *string
	indentLevel int
}

func (e *emitter) retarget(to *string) {
	e.target = to
}

func (e *emitter) writeln(f string, d ...interface{}) {
	e.write(f+"\n", d...)
}

func (e *emitter) writet(num int, f string, d ...interface{}) {
	const TabSize = 4
	tabs := strings.Repeat(" ", TabSize*num)
	e.write(tabs+f, d...)
}

func (e *emitter) writetln(num int, f string, d ...interface{}) {
	e.writet(num, f+"\n", d...)
}

func (e *emitter) write(f string, d ...interface{}) {
	*e.target += fmt.Sprintf(f, d...)
}

func (e *emitter) writeType(typ ir.Type) string {
	switch t := typ.(type) {

	// compile to uint8_t, uint16_t, etc.
	case *ir.IntegerType:
		res := fmt.Sprintf("int%d_t", t.Width)
		if !t.Signed {
			res = "u" + res
		}
		return res
	case *ir.FloatingType:
		if t.Width == 32 {
			return "float"
		}
		return "double"

	case *ir.VoidType:
		return "void"

	case *ir.ReferenceType:
		// FIXME
		return t.Name

	case *ir.PointerType:
		return fmt.Sprintf("%s*", e.writeType(t.Base))

	default:
		panic(fmt.Sprintf("unhandled type %s", reflect.TypeOf(t)))
	}

}

func (e *emitter) buildAlloca(a *ir.Alloca) {
	aType := e.writeType(a.Type)
	e.writetln(e.indentLevel, "%s %s = malloc(sizeof(*%s));", aType, a.Name, a.Name)
}

func (e *emitter) buildBuiltin(b *ir.Builtin) string {
	bType := e.writeType(b.Type)
	switch b.Name {
	case "sizeof":
		return fmt.Sprintf("sizeof(%s)", bType)
	case "make":
		return fmt.Sprintf("malloc(sizeof(%s))", bType)
	default:
		panic(fmt.Sprintf("unimplemented builtin %s", b.Name))
	}
}

func (e *emitter) buildExpr(l ir.Value) string {
	switch val := l.(type) {
	case *ir.IntegerValue:
		return val.RawValue.String()

	case *ir.FloatingValue:
		return fmt.Sprintf("%f", val.Value)

	case *ir.StringValue:
		return val.Value

	case *ir.BinaryExpression:
		lh := e.buildExpr(val.LHand)
		rh := e.buildExpr(val.RHand)
		return fmt.Sprintf("(%s%s%s)", lh, val.Op, rh)

	case *ir.Grouping:
		value := e.buildExpr(val.Val)
		return fmt.Sprintf("(%s)", value)

	case *ir.Identifier:
		return val.Name.Value

	case *ir.Builtin:
		return e.buildBuiltin(val)

	case *ir.UnaryExpression:
		value := e.buildExpr(val.Val)
		return fmt.Sprintf("(*%s)", value)

	case *ir.Assign:
		lh := e.buildExpr(val.LHand)
		rh := e.buildExpr(val.RHand)
		return fmt.Sprintf("%s %s %s", lh, val.Op, rh)

	case *ir.Call:
		lh := e.buildExpr(val.Left)
		var argsList string
		for idx, arg := range val.Params {
			if idx != 0 {
				argsList += ","
			}
			argsList += e.buildExpr(arg)
		}
		return fmt.Sprintf("%s(%s)", lh, argsList)

	default:
		panic(fmt.Sprintf("unimplemented expr %s", reflect.TypeOf(val)))
	}
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

	e.writetln(e.indentLevel, "%s%s %s%s", modifier, aType, l.Name, valueCode)
}

func (e *emitter) buildRet(r *ir.Return) {
	res := ";"
	if r.Val != nil {
		res = fmt.Sprintf(" %s;", e.buildExpr(r.Val))
	}
	e.writetln(e.indentLevel, "return%s", res)
}

func (e *emitter) buildLoop(l *ir.Loop) {
	e.writetln(e.indentLevel, "for(;;)")
	e.buildBlock(l.Body)
}

func (e *emitter) buildWhileLoop(w *ir.WhileLoop) {
	cond := e.buildExpr(w.Cond)
	var post string
	if w.Post != nil {
		post = e.buildExpr(w.Post)
	}
	e.writetln(e.indentLevel, "for(;%s;%s)", cond, post)
	e.buildBlock(w.Body)
}

func (e *emitter) buildIfStat(iff *ir.IfStatement) {
	e.writetln(e.indentLevel, "if(%s)", e.buildExpr(iff.Cond))
	e.buildBlock(iff.True)

	for _, elif := range iff.ElseIf {
		e.writetln(e.indentLevel, "else if(%s)", e.buildExpr(elif.Cond))
		e.buildBlock(elif.Body)
	}

	if iff.Else != nil {
		e.buildBlock(iff.Else)
	}
}

func (e *emitter) buildAssign(a *ir.Assign) {
	lh := e.buildExpr(a.LHand)
	op := a.Op
	rh := e.buildExpr(a.RHand)
	e.writetln(e.indentLevel, "%s %s %s;", lh, op, rh)
}

func (e *emitter) buildCall(c *ir.Call) {
	var argList string
	for idx, p := range c.Params {
		if idx != 0 {
			argList += ","
		}
		argList += e.buildExpr(p)
	}
	left := e.buildExpr(c.Left)

	// FIXME hard coded mangle thing
	e.writet(e.indentLevel, "%s(%s)", left, argList)
}

func (e *emitter) buildInstr(i ir.Instruction) {
	switch instr := i.(type) {

	// memory allocation
	case *ir.Alloca:
		e.buildAlloca(instr)
		return
	case *ir.Local:
		e.buildLocal(instr)
		return
	case *ir.Assign:
		e.buildAssign(instr)
		return

	// ...

	case *ir.Block:
		e.buildBlock(instr)
		return

	// ...

	case *ir.Return:
		e.buildRet(instr)
		return

	// conditional

	case *ir.IfStatement:
		e.buildIfStat(instr)
		return

	// looping constructs

	case *ir.Loop:
		e.buildLoop(instr)
		return

	case *ir.WhileLoop:
		e.buildWhileLoop(instr)
		return

	case *ir.Call:
		e.buildCall(instr)
		e.write(";\n")
		return

	default:
		panic(fmt.Sprintf("unhandled instr %s", reflect.TypeOf(instr)))
	}

}

func (e *emitter) buildBlock(b *ir.Block) {
	e.writetln(e.indentLevel, "{")
	e.indentLevel++

	for _, instr := range b.Instr {
		e.buildInstr(instr)
	}

	e.indentLevel--
	e.writetln(e.indentLevel, "}")
}

func (e *emitter) emitStructure(st *ir.Structure) {
	// forward declare 'Struct name' as just 'name'
	e.writetln(e.indentLevel, "typedef struct %s %s;", st.Name, st.Name)

	e.writetln(e.indentLevel, "struct %s {", st.Name)
	e.indentLevel++

	for _, name := range st.Fields.Order {
		t := st.Fields.Get(name.Value)

		typ := e.writeType(t)
		e.writetln(e.indentLevel, "%s %s;", typ, name)
	}
	e.indentLevel--

	e.writetln(e.indentLevel, "};")
}

func (e *emitter) emitFunc(fn *ir.Function) {
	mangledFuncName := fn.Name.Value
	if strings.Compare(fn.Name.Value, "main") == 0 {
		mangledFuncName = "krug_" + fn.Name.Value
	}

	writeArgList := func(fn *ir.Function) string {
		var argList string

		idx := 0
		for _, name := range fn.Param.Order {
			t := fn.Param.Get(name.Value)

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
	e.writeln("%s %s(%s);", typ, mangledFuncName, argList)

	// write definition to the source part.
	e.retarget(&e.source)
	e.writeln("%s %s(%s)", typ, mangledFuncName, argList)
	e.buildBlock(fn.Body)
}

func codegen(mod *ir.Module) ([]byte, []api.CompilerError) {
	e := &emitter{}
	e.retarget(&e.decl)

	headers := []string{
		"stdio.h",
		"stdbool.h",
		"stdint.h",
		"stdlib.h",
	}
	for _, h := range headers {
		e.writeln(`#include <%s>`, h)
	}

	for _, st := range mod.Structures {
		e.emitStructure(st)
	}

	e.retarget(&e.source)

	for _, fn := range mod.Functions {
		e.emitFunc(fn)
	}

	// for now we manually write the main func
	e.retarget(&e.source)
	e.writeln(`int main() { return krug_main(); }`)

	return []byte(e.decl + e.source), []api.CompilerError{}
}
