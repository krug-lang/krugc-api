package ir

import (
	"encoding/gob"
	"fmt"
)

type Type interface {
	String() string
}

func init() {
	gob.Register(&FloatingType{})
	gob.Register(&IntegerType{})
	gob.Register(&VoidType{})
	gob.Register(&ReferenceType{})
	gob.Register(&PointerType{})
	gob.Register(&Structure{})
	gob.Register(&Function{})
	gob.Register(&TypeDict{})
	gob.Register(&Impl{})
	gob.Register(&ArrayType{})
}

var (
	Float64 = NewFloatingType(64)
	Float32 = NewFloatingType(32)

	Int8  = NewIntegerType(8, true)
	Int16 = NewIntegerType(16, true)
	Int32 = NewIntegerType(32, true)
	Int64 = NewIntegerType(64, true)

	Uint8  = NewIntegerType(8, false)
	Uint16 = NewIntegerType(16, false)
	Uint32 = NewIntegerType(32, false)
	Uint64 = NewIntegerType(64, false)

	Void = &VoidType{}

	Bool = Uint32
	Rune = Int32
)

var PrimitiveType = map[string]Type{
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
	Width  int
	Signed bool
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
	Width int
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
	Name string
}

func (r *ReferenceType) String() string {
	return fmt.Sprintf("#%s", r.Name)
}

func NewReferenceType(name string) *ReferenceType {
	return &ReferenceType{name}
}

// ARRAY TYPE

type ArrayType struct {
	Base Type
	Size Value
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("[%s; %s]", a.Base.String(), a.Size)
}

func NewArrayType(base Type, size Value) *ArrayType {
	return &ArrayType{base, size}
}

// POINTER TYPE

type PointerType struct {
	Base Type
}

func (p *PointerType) String() string {
	return fmt.Sprintf("*%s", p.Base.String())
}

func NewPointerType(base Type) *PointerType {
	return &PointerType{base}
}

// ORDERED TYPE MAP.

type TypeDict struct {
	Data  map[string]Type
	Order []string
}

func (t *TypeDict) Get(k string) Type {
	typ, _ := t.Data[k]
	return typ
}

func (t *TypeDict) Set(a string, typ Type) {
	t.Order = append(t.Order, a)
	t.Data[a] = typ
}

func (t TypeDict) String() string {
	var fields string
	for _, name := range t.Order {
		field, _ := t.Data[name]
		fields += fmt.Sprintf("%s:%s,", name, field.String())
	}
	return fields
}

func newTypeDict() *TypeDict {
	return &TypeDict{map[string]Type{}, []string{}}
}

// STRUCTURE

type Structure struct {
	Name    string
	Fields  *TypeDict
	Methods map[string]*Function
}

func (s *Structure) RegisterMethod(f *Function) {
	s.Methods[f.Name] = f
}

func (s *Structure) String() string {
	return fmt.Sprintf("%s{%s}", s.Name, s.Fields)
}

func NewStructure(name string, fields *TypeDict) *Structure {
	return &Structure{name, fields, map[string]*Function{}}
}

// FUNCTION

type Function struct {
	Name       string
	Param      *TypeDict
	ReturnType Type
	Body       *Block
}

func NewFunction(name string, params *TypeDict, ret Type) *Function {
	return &Function{name, params, ret, NewBlock()}
}

type UnclaimedMethod struct {
	Parent string
	Method *Function
}

// IMPL - method group

type Impl struct {
	Name    string
	Methods map[string]*Function
}

func (i *Impl) RegisterMethod(fn *Function) {
	i.Methods[fn.Name] = fn
}

func NewImpl(name string) *Impl {
	return &Impl{name, map[string]*Function{}}
}
