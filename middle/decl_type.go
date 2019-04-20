package middle

import (
	"reflect"

	"github.com/hugobrains/caasper/api"
	"github.com/hugobrains/caasper/ir"
)

type decl struct {
	mod      *ir.Module          `json:"mod"`
	scopeMap *ir.ScopeMap        `json:"scope_map"`
	errors   []api.CompilerError `json:"errors"`
	curr     *ir.SymbolTable     `json:"curr"`
}

func (d *decl) error(e api.CompilerError) {
	d.errors = append(d.errors, e)
}

func (d *decl) push(stab *ir.SymbolTable) {
	d.curr = stab
}

func (d *decl) pop() {
	d.curr = d.curr.Outer
}

func (d *decl) regType(name string, t *ir.Type) {
	d.curr.RegisterType(name, t)
}

func (d *decl) visitLocal(l *ir.Local) {
	if l.Type == nil {
		d.error(api.NewUnimplementedError("type inference"))
		return
	}

	// if the type its a reference type,
	// try and link this to the type it references.
	if l.Type.Kind == ir.ReferenceKind {
		ref := l.Type.Reference
		name := ref.Name

		if structure, ok := d.mod.Structures[name]; ok {
			l.Type = &ir.Type{
				Kind:      ir.StructKind,
				Structure: structure,
			}
		} else {
			// couldn't find reference type.
		}
	}

	d.regType(l.Name.Value, l.Type)
}

func (d *decl) visitAlloca(a *ir.Alloca) {
	if a.Type == nil {
		d.error(api.NewUnimplementedError("type inference"))
		return
	}
	d.regType(a.Name.Value, a.Type)
}

func (d *decl) visitInstr(i *ir.Instruction) {
	switch i.Kind {
	case ir.BlockInstr:
		d.visitBlock(i.Block)

	case ir.LocalInstr:
		d.visitLocal(i.Local)
	case ir.AllocaInstr:
		d.visitAlloca(i.Alloca)

		// FIXME?
	case ir.PathValue:
		return

	case ir.ReturnInstr:
		return

	default:
		d.error(api.NewUnimplementedError("visitInstr: " + reflect.TypeOf(i).String()))
	}
}

func (d *decl) visitBlock(b *ir.Block) {
	// TODO restore the stab from the scope map
	for _, instr := range b.Instr {
		d.visitInstr(instr)
	}
	// POP
}

func declType(scopeMap *ir.ScopeMap, mod *ir.Module) (*ir.TypeMap, []api.CompilerError) {
	d := &decl{
		mod,
		scopeMap,
		[]api.CompilerError{},
		nil,
	}

	tm := ir.NewTypeMap()

	for _, name := range mod.FunctionOrder {
		fn, _ := mod.Functions[name.Value]

		fnStab := scopeMap.Functions[name.Value]
		d.push(fnStab)

		for _, instr := range fn.Body.Instr {
			d.visitInstr(instr)
		}

		d.pop()
	}

	return tm, d.errors
}
