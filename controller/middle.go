package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/caasper/middle"
	"github.com/krug-lang/caasper/service"
)

type middleService struct {}

func (ms *middleService) register(router *gin.Engine) {
	// 'middle' of compiler stages, takes the
	// krug IR and type checks everything.
	m := router.Group("/mid")
	{
		b := m.Group("/build")
		{
			// module -> [build/scope] -> scope map.
			b.POST("/scope", service.BuildScope)

			// module -> [build/scope_map] -> scope dict
			b.POST("/scope_dict", service.BuildScopeDict)

			// module -> [build/type] -> type map.
			b.POST("/type", service.BuildType)
		}

		m.POST("/unused_func", service.UnusedFunctions)

		// TODO grouping for these?
		m.POST("/borrow_check", middle.BorrowCheck)
		m.POST("mut_check", middle.MutabilityCheck)

		r := m.Group("/resolve")
		{
			// semantic module{module, scope_map, type_map} -> [type_resolve]
			//
			// this route takes in a module, scope map, and type map
			// and will resolve all of the types in top level declarations
			// e.g.
			// struct Person { name int, foo SomeStruct }
			// will check that all of the types in the struct resolve.
			r.POST("/type", middle.TypeResolve)

			// semantic module{module, scope_map, type_map} -> [symbol_resolve]
			//
			// this route takes in a moudle, scope map, and type map
			// it will check that all of the symbols exist in expressions
			r.POST("/symbol", middle.SymbolResolve)
		}
	}
}
