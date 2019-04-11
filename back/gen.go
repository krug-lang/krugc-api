package back

import (
	"fmt"
	"strings"

	"github.com/hugobrains/caasper/front"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
)

func Gen(c *gin.Context) {
	var codeGenReq api.CodeGenerationRequest
	if err := c.BindJSON(&codeGenReq); err != nil {
		panic(err)
	}

	var irMod ir.Module
	if err := jsoniter.Unmarshal([]byte(codeGenReq.IRModule), &irMod); err != nil {
		panic(err)
	}

	// for now we just return the
	// bytes for one big old c file.
	monoFile, errors := codegen(&irMod, codeGenReq.TabSize, codeGenReq.Minify)

	resp := api.KrugResponse{
		Data:   monoFile,
		Errors: errors,
	}
	c.JSON(200, &resp)
}

type emitter struct {
	decl   string
	source string
	target *string

	// this is a state representing
	// how many levels of indentation
	// the emitter is currently in.
	indentLevel int

	// this is the size of a tab in spaces
	tabSize int

	// whether or not the output is minified.
	minify bool
}

func (e *emitter) retarget(to *string) {
	e.target = to
}

// writeln will emitt the given line to the generated code
// NOTE that even if the minification is enabled, this will
// always generate a newline
func (e *emitter) writeln(f string, d ...interface{}) {
	e.write(f+"\n", d...)
}

// writet will emit a tabbed string to the target
func (e *emitter) writet(num int, f string, d ...interface{}) {
	tabs := " "
	if !e.minify {
		tabs = strings.Repeat(" ", e.tabSize*num)
	}

	e.write(tabs+f, d...)
}

// writetln 'write tabbed line' will write a tabbed (indented)
// line to the target.
// NOTE: if minification is enabled, new lines will not be appended here.
func (e *emitter) writetln(num int, f string, d ...interface{}) {
	newLine := "\n"
	if e.minify {
		newLine = ""
	}

	e.writet(num, f+newLine, d...)
}

// write will write a string to the given target.
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

func (e *emitter) emitTupleType(tup *ir.TupleType) string {
	var structTypes string
	for idx, typ := range tup.Types {
		if idx != 0 {
			structTypes += " "
		}
		structTypes += fmt.Sprintf("%s _%d;", e.writeType(typ), idx)
	}
	return fmt.Sprintf("struct { %s }", structTypes)
}

// HM
func (e *emitter) writeArray(typ *ir.ArrayType) string {
	return fmt.Sprintf("%s*", e.writeType(typ.Base))
}

func (e *emitter) writeType(typ *ir.Type) string {
	if typ == nil {
		panic("trying to write a nil type!")
	}

	switch typ.Kind {

	// compile to uint8_t, uint16_t, etc.
	case ir.IntegerKind:
		return e.writeInteger(typ.IntegerType)
	case ir.FloatKind:
		return e.writeFloat(typ.FloatingType)
	case ir.VoidKind:
		return "void"

	case ir.TupleKind:
		return e.emitTupleType(typ.Tuple)

	case ir.ReferenceKind:
		return typ.Reference.Name

	case ir.ArrayKind:
		return e.writeArray(typ.ArrayType)
	case ir.PointerKind:
		return e.writePointer(typ.Pointer)

	default:
		panic(fmt.Sprintf("unhandled type %s", typ.Kind))
	}

}

func (e *emitter) buildAlloca(a *ir.Alloca) {
	aType := e.writeType(a.Type)
	e.writetln(e.indentLevel, "%s %s = malloc(sizeof(*%s));", aType, a.Name, a.Name)

	// TODO init list.
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

func (e *emitter) writeInitExpr(i *ir.Init) string {
	switch i.Kind {
	case front.InitStructure:
		return ""
	case front.InitTuple:
		return ""
	case front.InitArray:
		return ""
	}
	panic(fmt.Sprintf("unimplemented int expr %s", i.Kind))
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

	case ir.InitValue:
		return e.writeInitExpr(l.Init)

	default:
		panic(fmt.Sprintf("unimplemented expr %s", l.Kind))
	}
}

func (e *emitter) removePointer(ptr *ir.Type) *ir.Type {
	if ptr.Kind == ir.PointerKind {
		return ptr.Pointer.Base
	}
	return ptr
}

func (e *emitter) buildInitializerFor(l *ir.Local, init *ir.Init) {
	localName := l.Name.Value

	switch init.Kind {
	case front.InitTuple:
		for idx, val := range init.Values {
			genVal := e.buildExpr(val)
			e.writetln(e.indentLevel, "%s._%d = %s;", localName, idx, genVal)
		}

	case front.InitStructure:
		break

	case front.InitArray:
		var genLit string
		for idx, expr := range init.Values {
			if idx != 0 {
				genLit += ","
			}
			genLit += e.buildExpr(expr)
		}

		// this is a massive hack for now, all arrays are
		// changed into pointer types

		// we then create an actual array and set the pointer
		// to point at it.

		// HACK, all arrays are changed into a pointer type.
		genType := e.writeType(e.removePointer(l.Type))
		e.writetln(e.indentLevel, "%s _%s_raw_arr[] = {%s};", genType, localName, genLit)

		e.writetln(e.indentLevel, "%s = _%s_raw_arr;", localName, localName)
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
		// initializer is emitted AFTER the variable.
		if l.Val.Kind != ir.InitValue {
			valueCode = fmt.Sprintf(" = %s;", e.buildExpr(l.Val))
		}
	}

	e.writetln(e.indentLevel, "%s%s %s%s", modifier, aType, l.Name.Value, valueCode)

	if l.Val != nil && l.Val.Kind == ir.InitValue {
		e.buildInitializerFor(l, l.Val.Init)
	}
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
		e.writetln(e.indentLevel, "else")
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

	case ir.ExpressionInstr:
		e.writetln(e.indentLevel, "%s;", e.buildExpr(i.ExpressionStatement))
		return

	default:
		panic(fmt.Sprintf("unhandled instr %s", i.Kind))
	}

}

func (e *emitter) buildBlock(b *ir.Block) {
	e.writetln(e.indentLevel, "{")
	e.indentLevel++

	for _, instr := range b.Instr {
		e.buildInstr(instr)
	}

	for i := len(b.DeferStack) - 1; i >= 0; i-- {
		item := b.DeferStack[i]
		if item.Block != nil {
			e.buildBlock(item.Block)
		} else {
			e.buildInstr(item.Stat)
		}
	}

	if b.Return != nil {
		e.buildInstr(b.Return)
	}

	e.indentLevel--
	e.writetln(e.indentLevel, "}")
}

func (e *emitter) emitStructure(st *ir.Structure) {
	stName := st.Name.Value

	// forward declare 'Struct name' as just 'name'
	e.writetln(e.indentLevel, "typedef struct %s %s;", stName, stName)

	e.writetln(e.indentLevel, "struct %s {", stName)
	e.indentLevel++

	for _, name := range st.Fields.Order {
		t := st.Fields.Get(name.Value)

		typ := e.writeType(t)
		e.writetln(e.indentLevel, "%s %s;", typ, name.Value)
	}
	e.indentLevel--

	e.writetln(e.indentLevel, "};")
}

func (e *emitter) emitFunc(fn *ir.Function) {
	generatedFuncName := fn.Name.Value

	// when the function is just 'main', we mangle it
	// to krug_main, as this is the entry point of our program.
	if strings.Compare(fn.Name.Value, "main") == 0 {
		generatedFuncName = "krug_" + fn.Name.Value
	}

	writeArgList := func(fn *ir.Function) string {
		var argList string

		idx := 0
		for _, name := range fn.Param.Order {
			t := fn.Param.Get(name.Value)

			if idx != 0 {
				argList += ", "
			}
			argList += fmt.Sprintf("%s %s", e.writeType(t), name.Value)
			idx++
		}
		return argList
	}

	returnType := e.writeType(fn.ReturnType)
	argList := writeArgList(fn)

	// write prototype to the decl part
	e.retarget(&e.decl)
	e.writeln("%s %s(%s);", returnType, generatedFuncName, argList)

	// write definition to the source part.
	e.retarget(&e.source)

	e.writeln("%s %s(%s)", returnType, generatedFuncName, argList)

	e.buildBlock(fn.Body)
}

// TODO if this arg list gets any bigger it should become
// some kind of configuration struct.
func codegen(mod *ir.Module, tabSize int, minify bool) (string, []api.CompilerError) {
	fmt.Println(mod)

	e := &emitter{
		tabSize: tabSize,
		minify:  minify,
	}
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
