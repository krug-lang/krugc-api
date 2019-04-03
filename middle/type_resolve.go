package middle

import (
	"fmt"
	"reflect"

	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
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

func (t *typeResolvePass) resolveReferenceType(ref *ir.ReferenceType) *ir.Type {
	// resolve the type:

	// 1. primitive type?
	// this check would have been done during buildType
	// so this is superfluous but ok
	if typ, ok := ir.PrimitiveType[ref.Name]; ok {
		return typ
	}

	// 2. structure type?
	if typ, ok := t.mod.Structures[ref.Name]; ok {
		return &ir.Type{
			Kind:      ir.StructKind,
			Structure: typ,
		}
	}

	// 3. trait type?

	t.error(api.CompilerError{
		Title: fmt.Sprintf("Couldn't resolve type '%s'", ref.Name),
		Desc:  "",
	})
	return nil
}

func (t *typeResolvePass) resolveVia(typ *ir.Type, val *ir.Value) *ir.Type {
	if typ == nil {

		if val.Kind != ir.IdentifierValue {
			panic(fmt.Sprintf("left is %s", reflect.TypeOf(val).String()))
		}

		scope := t.scope()
		if scope == nil {
			panic("what zeee fuck?")
		}

		iden := val.Identifier

		typ, ok := scope.LookupType(iden.Name.Value)
		if !ok {
			t.error(api.NewUnresolvedSymbol(iden.Name.Value, iden.Name.Span...))
		}

		return typ
	}

	switch typ.Kind {
	case ir.StructKind:
		iden := val.Identifier

		left := typ.Structure
		typ, ok := left.Fields.Data[iden.Name.Value]
		if !ok {
			t.error(api.NewUnresolvedSymbol(iden.Name.Value, iden.Name.Span...))
		}
		return typ

	default:
		t.error(api.NewUnimplementedError(reflect.TypeOf(typ).String()))
	}

	return nil
}

func (t *typeResolvePass) resolvePath(p *ir.Path) {
	var lastType *ir.Type
	for _, val := range p.Values {
		lastType = t.resolveVia(lastType, val)
	}
}

func (t *typeResolvePass) resolveType(typ *ir.Type) *ir.Type {
	switch typ.Kind {
	case ir.PointerKind:
		return t.resolveType(typ.Pointer.Base)

	case ir.ReferenceKind:
		return t.resolveReferenceType(typ.Reference)

	case ir.ArrayKind:
		return t.resolveType(typ.ArrayType.Base)

	case ir.StructKind:
		return t.resolveStructure(typ.Structure)

	case ir.IntegerKind:
		return typ
	case ir.FloatKind:
		return typ
	case ir.VoidKind:
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

func (t *typeResolvePass) resolveInstr(i *ir.Instruction) {
	switch i.Kind {
	case ir.AllocaInstr:
		t.resolveAlloca(i.Alloca)
	case ir.LocalInstr:
		t.resolveLocal(i.Local)
	case ir.WhileLoopInstr:
		t.resolveBlock(i.WhileLoop.Body)
	case ir.LoopInstr:
		t.resolveBlock(i.Loop.Body)

	case ir.IfStatementInstr:
		instr := i.IfStatement
		t.resolveBlock(instr.True)
		for _, e := range instr.ElseIf {
			t.resolveBlock(e.Body)
		}
		if instr.Else != nil {
			t.resolveBlock(instr.Else)
		}

	case ir.PathValue:
		// i think this may have to go into
		// a later pass.
		// first pass resolves top level decls
		// then the second pass resolves the
		// func level types?
		// t.resolvePath(instr)

	case ir.ReturnInstr:
		return
	case ir.BreakInstr:
		return
	case ir.NextInstr:
		return

	// FIXME
	case ir.CallValue:
		return

	case ir.AssignInstr:
		return

	case ir.BlockInstr:
		t.resolveBlock(i.Block)

	default:
		panic(fmt.Sprintf("unhandled instruction %s", reflect.TypeOf(i)))
	}
}

func (t *typeResolvePass) resolveBlock(b *ir.Block) {
	t.push(b.Stab)
	for _, instr := range b.Instr {
		t.resolveInstr(instr)
	}
	t.pop()
}

func (t *typeResolvePass) resolveStructure(st *ir.Structure) *ir.Type {
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

	// HM
	return &ir.Type{
		Kind:      ir.StructKind,
		Structure: st,
	}
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
