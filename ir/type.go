package ir

import (
	"fmt"

	"github.com/krug-lang/caasper/front"
)

type TypeOf string

const (
	FloatKind     TypeOf = "float"
	IntegerKind          = "int"
	ArrayKind            = "array"
	FunctionKind         = "fn"
	VoidKind             = "void"
	StructKind           = "struct"
	PointerKind          = "ptr"
	TupleKind            = "tuple"
	ReferenceKind        = "ref"
)

type Type struct {
	Kind TypeOf `json:"kind"`

	VoidType     *VoidType      `json:"voidType,omitempty"`
	FloatingType *FloatingType  `json:"floatingType,omitempty"`
	IntegerType  *IntegerType   `json:"integerType,omitempty"`
	ArrayType    *ArrayType     `json:"arrayType,omitempty"`
	Function     *Function      `json:"function,omitempty"`
	Tuple        *TupleType     `json:"tuple,omitempty"`
	Structure    *Structure     `json:"structure,omitempty"`
	Pointer      *PointerType   `json:"pointer,omitempty"`
	Reference    *ReferenceType `json:"reference,omitempty"`
}

func (t *Type) String() string {
	switch t.Kind {
	case VoidKind:
		return t.VoidType.String()
	case FloatKind:
		return t.FloatingType.String()
	case IntegerKind:
		return t.IntegerType.String()
	case ArrayKind:
		return t.ArrayType.String()
	case FunctionKind:
		return t.Function.String()
	case StructKind:
		return t.Structure.String()
	case PointerKind:
		return t.Pointer.String()
	case ReferenceKind:
		return t.Reference.String()
	default:
		panic("unhandled type in Type::String()")
	}
}

var (
	Float64 = &Type{Kind: FloatKind, FloatingType: NewFloatingType(64)}
	Float32 = &Type{Kind: FloatKind, FloatingType: NewFloatingType(32)}

	Int8  = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(8, true)}
	Int16 = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(16, true)}
	Int32 = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(32, true)}
	Int64 = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(64, true)}

	Uint8  = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(8, false)}
	Uint16 = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(16, false)}
	Uint32 = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(32, false)}
	Uint64 = &Type{Kind: IntegerKind, IntegerType: NewIntegerType(64, false)}

	Void = &Type{Kind: VoidKind, VoidType: &VoidType{}}

	Bool = Uint32
	Rune = Int32
)

var PrimitiveType = map[string]*Type{
	"f64": Float64,
	"f32": Float32,

	"i8":  Int8,
	"i16": Int16,
	"i32": Int32,
	"i64": Int64,

	"u8":  Uint8,
	"u16": Uint16,
	"u32": Uint32,
	"u64": Uint64,

	"void": Void,
	"bool": Bool,
	"rune": Rune,

	"int":  Int32,
	"uint": Uint32,
}

type VoidType struct{}

func (v *VoidType) String() string {
	return "void"
}

// INTEGER TYPE

type IntegerType struct {
	Width  int  `json:"width"`
	Signed bool `json:"signed"`
}

func (i *IntegerType) String() string {
	sign := "u"
	if i.Signed {
		sign = "s"
	}
	return fmt.Sprintf("%sint%d", sign, i.Width)
}

func NewIntegerType(width int, signed bool) *IntegerType {
	return &IntegerType{width, signed}
}

// FLOATING TYPE

type FloatingType struct {
	Width int `json:"width"`
}

func (f *FloatingType) String() string {
	return fmt.Sprintf("flt%d", f.Width)
}

func NewFloatingType(width int) *FloatingType {
	return &FloatingType{width}
}

// REFERENCE TYPE

// i.e. a reference to some
// type, e.g. Person or Shape, etc.
type ReferenceType struct {
	Name string `json:"name"`
}

func (r *ReferenceType) String() string {
	return fmt.Sprintf("#%s", r.Name)
}

func NewReferenceType(name string) *ReferenceType {
	return &ReferenceType{name}
}

// TUPLE TYPE

type TupleType struct {
	Types []*Type `json:"types"`
}

func (t *TupleType) String() string {
	return fmt.Sprintf("%s", t.Types)
}

func NewTupleType(types []*Type) *TupleType {
	return &TupleType{types}
}

// ARRAY TYPE

type ArrayType struct {
	Base *Type  `json:"base"`
	Size *Value `json:"size"`
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("[%s; %s]", a.Base.String(), a.Size)
}

func NewArrayType(base *Type, size *Value) *ArrayType {
	return &ArrayType{base, size}
}

// POINTER TYPE

type PointerType struct {
	Base *Type `json:"base"`
}

func (p *PointerType) String() string {
	return fmt.Sprintf("*%s", p.Base.String())
}

func NewPointerType(base *Type) *PointerType {
	return &PointerType{base}
}

// ORDERED TYPE MAP.

// TODO rename this to LocalDict
// maybe even rename a Local to a Binding
// since we bind a name to a type.
type TypeDict struct {
	Data  map[string]*Local `json:"data"`
	Order []front.Token     `json:"order"`
}

func (t *TypeDict) Get(name string) *Local {
	loc, _ := t.Data[name]
	return loc
}

func (t *TypeDict) Add(l *Local) {
	t.Order = append(t.Order, l.Name)
	t.Data[l.Name.Value] = l
}

func (t TypeDict) String() string {
	var fields string
	for _, name := range t.Order {
		field, _ := t.Data[name.Value]
		// FIXME?
		fields += fmt.Sprintf("%s:%s,", name, field)
	}
	return fields
}

func newTypeDict() *TypeDict {
	return &TypeDict{map[string]*Local{}, []front.Token{}}
}

// STRUCTURE

type Structure struct {
	Name    front.Token          `json:"name"`
	Stab    *SymbolTable         `json:"stab,omitempty"`
	Fields  *TypeDict            `json:"fields"`
	Methods map[string]*Function `json:"methods,omitempty"`
}

func (s *Structure) RegisterMethod(f *Function) {
	s.Methods[f.Name.Value] = f
}

func (s *Structure) String() string {
	return fmt.Sprintf("%s{%s}", s.Name, s.Fields)
}

func NewStructure(name front.Token, fields *TypeDict) *Structure {
	return &Structure{name, nil, fields, map[string]*Function{}}
}

// FUNCTION

type Function struct {
	Name       front.Token  `json:"name"`
	Stab       *SymbolTable `json:"stab,omitempty"`
	Param      *TypeDict    `json:"param"`
	ReturnType *Type        `json:"return_type,omitempty"`
	Body       *Block       `json:"body"`
}

func (f *Function) String() string {
	return f.ReturnType.String()
}

func NewFunction(name front.Token, params *TypeDict, ret *Type) *Function {
	return &Function{name, nil, params, ret, NewBlock()}
}

type UnclaimedMethod struct {
	Parent string    `json:"parent"`
	Method *Function `json:"method"`
}

// IMPL - method group

type Impl struct {
	Name    front.Token          `json:"name"`
	Stab    *SymbolTable         `json:"stab,omitempty"`
	Methods map[string]*Function `json:"methods"`
}

func (i *Impl) RegisterMethod(fn *Function) bool {
	if _, ok := i.Methods[fn.Name.Value]; ok {
		return false
	}
	i.Methods[fn.Name.Value] = fn
	return true
}

func NewImpl(name front.Token) *Impl {
	return &Impl{name, nil, map[string]*Function{}}
}
