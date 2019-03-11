package middle

import (
	"github.com/krug-lang/krugc-api/api"
	"github.com/krug-lang/krugc-api/ir"
)

type symResolvePass struct {
	mod    *ir.Module
	errors []api.CompilerError
}

func (s *symResolvePass) error(err api.CompilerError) {
	s.errors = append(s.errors, err)
}

func symResolve(mod *ir.Module) []api.CompilerError {
	srp := &symResolvePass{mod, []api.CompilerError{}}
	return srp.errors
}
