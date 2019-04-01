package middle

import (
	"fmt"
	"reflect"

	"github.com/hugobrains/krug-serv/api"
	"github.com/hugobrains/krug-serv/ir"
)

type typeResolvePass struct {
	mod    *ir.Module
	errors []api.CompilerError
	curr   *ir.SymbolTable
}

func (t *typeResolvePass) scope() *ir.SymbolTable {
	return t.curr
}

func (t *typeResolvePass) push(s *ir.SymbolTable) {
	t.curr = s
}

func (t *typeResolvePass) pop() {
	t.curr = t.curr.Outer
}

func (t *typeResolvePass) error(err api.CompilerError) {
	t.errors = append(t.errors, err)
}

func (t *typeResolvePass) resolveReferenceType(ref *ir.ReferenceType) ir.Type {
	// resolve the type:

	// 1. primitive type?
	// this check would have been done during buildType
	// so this is superfluous but ok
	if typ, ok := ir.PrimitiveType[ref.Name]; ok {
		return typ
	}

	// 2. structure type?
	if typ, ok := t.mod.Structures[ref.Name]; ok {
		return typ
	}

	// 3. trait type?

	t.error(api.CompilerError{
		Title: fmt.Sprintf("Couldn't resolve type '%s'", ref.Name),
		Desc:  "",
	})
	return nil
}

func (t *typeResolvePass) resolveVia(typ ir.Type, val ir.Value) ir.Type {
	if typ == nil {
		// should be an identifier ?
		iden, ok := val.(*ir.Identifier)
		if !ok {
			panic(fmt.Sprintf("left is %s", reflect.TypeOf(val).String()))
		}

		scope := t.scope()
		if scope == nil {
			panic("what zeee fuck?")
		}

		typ, ok := scope.LookupType(iden.Name.Value)
		if !ok {
			t.error(api.NewUnresolvedSymbol(iden.Name.Value, iden.Name.Span...))
		}
		return typ
	}

	switch left := typ.(type) {
	case *ir.Structure:
		iden, _ := val.(*ir.Identifier)
		typ, ok := left.Fields.Data[iden.Name.Value]
		if !ok {
			t.error(api.NewUnresolvedSymbol(iden.Name.Value, iden.Name.Span...))
		}
		return typ

	default:
		t.error(api.NewUnimplementedError(reflect.TypeOf(left).String()))
	}

	return nil
}

func (t *typeResolvePass) resolvePath(p *ir.Path) {
	var lastType ir.Type
	for _, val := range p.Values {
		lastType = t.resolveVia(lastType, val)
	}
}

func (t *typeResolvePass) resolveType(unresolvedType ir.Type) ir.Type {
	switch typ := unresolvedType.(type) {
	case *ir.PointerType:
		return t.resolveType(typ.Base)

	case *ir.ReferenceType:
		return t.resolveReferenceType(typ)

	case *ir.ArrayType:
		return t.resolveType(typ.Base)

	case *ir.Structure:
		return t.resolveStructure(typ)

	case *ir.IntegerType:
		return typ
	case *ir.FloatingType:
		return typ
	case *ir.VoidType:
		return typ

	default:
		panic(fmt.Sprintf("unhandled type %s", reflect.TypeOf(typ)))
		return nil
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

	case *ir.Path:
		// i think this may have to go into
		// a later pass.
		// first pass resolves top level decls
		// then the second pass resolves the
		// func level types?
		// t.resolvePath(instr)

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
	t.push(b.Stab)
	for _, instr := range b.Instr {
		t.resolveInstr(instr)
	}
	t.pop()
}

func (t *typeResolvePass) resolveStructure(st *ir.Structure) ir.Type {
	for _, fn := range st.Methods {
		t.resolveFunc(fn)
	}

	for _, name := range st.Fields.Order {
		typ, _ := st.Fields.Data[name.Value]

		// set the field type to the newly
		// resolved type
		resolved := t.resolveType(typ)
		st.Fields.Data[name.Value] = resolved
	}

	return st
}

func (t *typeResolvePass) resolveFunc(fn *ir.Function) {
	for _, name := range fn.Param.Order {
		param, _ := fn.Param.Data[name.Value]
		t.resolveType(param)
	}

	t.resolveType(fn.ReturnType)

	t.push(fn.Stab)
	for _, instr := range fn.Body.Instr {
		t.resolveInstr(instr)
	}
	t.pop()
}

func typeResolve(mod *ir.Module) (*ir.TypeMap, []api.CompilerError) {
	trp := &typeResolvePass{mod, []api.CompilerError{}, nil}

	// todo get this from the scope map
	// trp.push(mod.Root)

	tm := ir.NewTypeMap()

	for _, impl := range mod.Impls {
		structure, ok := mod.GetStructure(impl.Name.Value)
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

	return tm, trp.errors
}
