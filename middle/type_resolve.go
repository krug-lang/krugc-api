package middle

import (
	"fmt"
	"reflect"

	"github.com/krug-lang/krugc-api/api"
	"github.com/krug-lang/krugc-api/ir"
)

type typeResolvePass struct {
	mod    *ir.Module
	errors []api.CompilerError
}

func (t *typeResolvePass) error(err api.CompilerError) {
	t.errors = append(t.errors, err)
}

func (t *typeResolvePass) resolveReferenceType(ref *ir.ReferenceType) {
	// resolve the type:

	// 1. primitive type?
	// this check would have been done during buildType
	// so this is superfluous but ok
	if _, ok := ir.PrimitiveType[ref.Name]; ok {
		return
	}

	// 2. structure type?
	if _, ok := t.mod.Structures[ref.Name]; ok {
		return
	}

	// 3. trait type?

	t.error(api.CompilerError{
		Title: fmt.Sprintf("Couldn't resolve type '%s'", ref.Name),
		Desc:  "",
	})
}

func (t *typeResolvePass) resolveType(unresolvedType ir.Type) {
	switch typ := unresolvedType.(type) {
	case *ir.PointerType:
		t.resolveType(typ.Base)
	case *ir.ReferenceType:
		t.resolveReferenceType(typ)
	case *ir.ArrayType:
		t.resolveType(typ.Base)

	case *ir.IntegerType:
		return
	case *ir.FloatingType:
		return
	case *ir.VoidType:
		return

	default:
		panic(fmt.Sprintf("unhandled type %s", reflect.TypeOf(typ)))
	}
}

func (t *typeResolvePass) resolveAlloca(a *ir.Alloca) {
	t.resolveType(a.Type)
}

func (t *typeResolvePass) resolveLocal(l *ir.Local) {
	t.resolveType(l.Type)
}

func (t *typeResolvePass) resolveInstr(i ir.Instruction) {
	switch instr := i.(type) {
	case *ir.Alloca:
		t.resolveAlloca(instr)
	case *ir.Local:
		t.resolveLocal(instr)
	case *ir.WhileLoop:
		t.resolveBlock(instr.Body)
	case *ir.Loop:
		t.resolveBlock(instr.Body)

	case *ir.IfStatement:
		t.resolveBlock(instr.True)
		for _, e := range instr.ElseIf {
			t.resolveBlock(e.Body)
		}
		if instr.Else != nil {
			t.resolveBlock(instr.Else)
		}

	// nops
	case *ir.Path:
		return
	case *ir.Return:
		return
	case *ir.Break:
		return
	case *ir.Next:
		return
	case *ir.Call:
		return
	case *ir.Assign:
		return

	case *ir.Block:
		t.resolveBlock(instr)

	default:
		panic(fmt.Sprintf("unhandled instruction %s", reflect.TypeOf(instr)))
	}
}

func (t *typeResolvePass) resolveBlock(b *ir.Block) {
	for _, instr := range b.Instr {
		t.resolveInstr(instr)
	}
}

func (t *typeResolvePass) resolveStructure(st *ir.Structure) {
	for _, fn := range st.Methods {
		t.resolveFunc(fn)
	}

	for _, name := range st.Fields.Order {
		typ, _ := st.Fields.Data[name]
		t.resolveType(typ)
	}
}

func (t *typeResolvePass) resolveFunc(fn *ir.Function) {
	for _, name := range fn.Param.Order {
		param, _ := fn.Param.Data[name]
		t.resolveType(param)
	}

	t.resolveType(fn.ReturnType)
	t.resolveBlock(fn.Body)
}

func typeResolve(mod *ir.Module) []api.CompilerError {
	trp := &typeResolvePass{mod, []api.CompilerError{}}

	for _, impl := range mod.Impls {
		structure, ok := mod.GetStructure(impl.Name)
		if !ok {
			trp.error(api.CompilerError{
				Title: fmt.Sprintf("Couldn't resolve structure '%s' being implemented", impl.Name),
				Desc:  "...",
			})
		}

		for _, fn := range impl.Methods {
			trp.resolveFunc(fn)
			structure.RegisterMethod(fn)
		}
	}

	for _, st := range mod.Structures {
		trp.resolveStructure(st)
	}

	for _, fn := range mod.Functions {
		trp.resolveFunc(fn)
	}

	return trp.errors
}
