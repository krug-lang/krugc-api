package entity

// sema stuff

type BuildTypeMapRequest struct{}

type BuildScopeDictRequest struct {
	IRModule string `json:"ir_module"`
}

type BuildScopeMapRequest struct {
	IRModule string `json:"ir_module"`
}

// unused func

// UnusedFunctionRequest is a request that will check
// for all unused functions
type UnusedFunctionRequest struct {
	IRModule string `json:"ir_module"`
	ScopeMap string `json:"scope_map"`
}

// borrowcheck

type BorrowCheckRequest struct {
	IRModule string `json:"ir_module"`
	ScopeMap string `json:"scope_map"`
}

// mut check

type MutabilityCheckRequest struct {
	IRModule string `json:"ir_module"`
	ScopeMap string `json:"scope_map"`
}

// resolution stuff

type TypeResolveRequest struct{}
type SymbolResolveRequest struct{}