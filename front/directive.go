package front

type valueKind string

const (
	stringValue    valueKind = "str"
	integerValue             = "int"
	characterValue           = "val"
	floatingValue            = "float"
)

type value struct {
	kind  valueKind
	value Token
}

type directiveKind string

const (
	Include  directiveKind = "include"
	Link                   = "link"
	NoMangle               = "no_mangle"
	Align                  = "align"
	Packed                 = "packed"
	Clang                  = "clang"
)

// #{include(string)}
type IncludeDirective struct {
	Path string `json:"path"`
}

// #{link("/some/path")}
type LinkDirective struct {
	Flags []string `json:"flags"`
}

// #{align(integer)}
type AlignDirective struct {
	Alignment uint64 `json:"align"`
}

// #{clang}
type ClangDirective struct{}

// #{no_mangle}
type NoMangleDirective struct{}

// #{packed}
type PackedDirective struct{}

type Directive struct {
	Kind              directiveKind
	IncludeDirective  *IncludeDirective  `json:"includeDirective,omitempty"`
	LinkDirective     *LinkDirective     `json:"linkDirective,omitempty"`
	AlignDirective    *AlignDirective    `json:"alignDirective,omitempty"`
	NoMangleDirective *NoMangleDirective `json:"noMangleDirective,omitempty"`
	PackedDirective   *PackedDirective   `json:"packedDirective,omitempty"`
	ClangDirective    *ClangDirective    `json:"clangDirective,omitempty"`
}
