package ir

import "encoding/gob"

type Type interface{}

func init() {
	gob.Register(&FloatingType{})
	gob.Register(&IntegerType{})
	gob.Register(&VoidType{})
	gob.Register(&ReferenceType{})
	gob.Register(&PointerType{})
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

type VoidType struct{}

// INTEGER TYPE

type IntegerType struct {
	Width  int
	Signed bool
}

func NewIntegerType(width int, signed bool) *IntegerType {
	return &IntegerType{width, signed}
}

// FLOATING TYPE

type FloatingType struct {
	Width int
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

func NewReferenceType(name string) *ReferenceType {
	return &ReferenceType{name}
}

// POINTER TYPE

type PointerType struct {
	Base Type
}

func NewPointerType(base Type) *PointerType {
	return &PointerType{base}
}
