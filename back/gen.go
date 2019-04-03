package back

import (
	"fmt"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
)

func Gen(c *gin.Context) {
	var krugReq api.KrugRequest
	if err := c.BindJSON(&krugReq); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(krugReq.Data), &irMod); err != nil {
		panic(err)
	}

	// for now we just return the
	// bytes for one big old c file.
	monoFile, errors := codegen(&irMod)

	resp := api.KrugResponse{
		Data:   monoFile,
		Errors: errors,
	}
	c.JSON(200, &resp)
}

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

func (e *emitter) writeInteger(t *ir.IntegerType) string {
	res := fmt.Sprintf("int%d_t", t.Width)
	if !t.Signed {
		res = "u" + res
	}
	return res
}

func (e *emitter) writeFloat(t *ir.FloatingType) string {
	if t.Width == 32 {
		return "float"
	}
	return "double"
}

func (e *emitter) writePointer(t *ir.PointerType) string {
	return fmt.Sprintf("%s*", e.writeType(t.Base))
}

func (e *emitter) writeType(typ *ir.Type) string {
	switch typ.Kind {

	// compile to uint8_t, uint16_t, etc.
	case ir.IntegerKind:
		return e.writeInteger(typ.IntegerType)
	case ir.FloatKind:
		return e.writeFloat(typ.FloatingType)
	case ir.VoidKind:
		return "void"

	case ir.ReferenceKind:
		return typ.Reference.Name

	case ir.PointerKind:
		return e.writePointer(typ.Pointer)

	default:
		panic(fmt.Sprintf("unhandled type %s", reflect.TypeOf(typ)))
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

func (e *emitter) buildExpr(l *ir.Value) string {
	switch l.Kind {
	case ir.IntegerValueValue:
		val := l.IntegerValue
		return val.RawValue.String()

	case ir.FloatingValueValue:
		val := l.FloatingValue
		return fmt.Sprintf("%f", val.Value)

	case ir.StringValueValue:
		val := l.StringValue
		return val.Value

	case ir.BinaryExpressionValue:
		val := l.BinaryExpression
		lh := e.buildExpr(val.LHand)
		rh := e.buildExpr(val.RHand)
		return fmt.Sprintf("(%s%s%s)", lh, val.Op, rh)

	case ir.GroupingValue:
		val := l.Grouping
		value := e.buildExpr(val.Val)
		return fmt.Sprintf("(%s)", value)

	case ir.IdentifierValue:
		val := l.Identifier
		return val.Name.Value

	case ir.BuiltinValue:
		val := l.Builtin
		return e.buildBuiltin(val)

	case ir.UnaryExpressionValue:
		val := l.UnaryExpression
		value := e.buildExpr(val.Val)
		return fmt.Sprintf("(*%s)", value)

	case ir.AssignValue:
		val := l.Assign
		lh := e.buildExpr(val.LHand)
		rh := e.buildExpr(val.RHand)
		return fmt.Sprintf("%s %s %s", lh, val.Op, rh)

	case ir.CallValue:
		val := l.Call
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
		panic(fmt.Sprintf("unimplemented expr %s", reflect.TypeOf(l)))
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

func (e *emitter) buildInstr(i *ir.Instruction) {
	switch i.Kind {

	// memory allocation
	case ir.AllocaInstr:
		e.buildAlloca(i.Alloca)
		return
	case ir.LocalInstr:
		e.buildLocal(i.Local)
		return
	case ir.AssignInstr:
		e.buildAssign(i.Assign)
		return

	// ...

	case ir.BlockInstr:
		e.buildBlock(i.Block)
		return

	// ...

	case ir.ReturnInstr:
		e.buildRet(i.Return)
		return

	// conditional

	case ir.IfStatementInstr:
		e.buildIfStat(i.IfStatement)
		return

	// looping constructs

	case ir.LoopInstr:
		e.buildLoop(i.Loop)
		return

	case ir.WhileLoopInstr:
		e.buildWhileLoop(i.WhileLoop)
		return

		/* hm?
		case ir.CallValue:
			e.buildCall(i.Call)
			e.write(";\n")
			return
		*/

	default:
		panic(fmt.Sprintf("unhandled instr %s", reflect.TypeOf(i)))
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

func codegen(mod *ir.Module) (string, []api.CompilerError) {
	fmt.Println(mod)

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

	return string(e.decl + e.source), []api.CompilerError{}
}
