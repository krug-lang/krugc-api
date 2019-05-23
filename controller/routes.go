package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krug-lang/caasper/back"
	"github.com/krug-lang/caasper/front"
	"github.com/krug-lang/caasper/ir"
	"github.com/krug-lang/caasper/middle"
)

func RegisterRoutes(router *gin.Engine) {
	const VERSION = "0.0.1"
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"version": VERSION,
		})
	})

	// compiler frontend, handles lexing/parsing
	f := router.Group("/front")
	{
		// lexical analysis
		f.POST("/lex", front.Tokenize)

		// parsing.
		parse := f.Group("/parse")
		{
			parse.POST("/ast", front.Parse)
			parse.POST("/directive", front.DirectiveParser)
		}

		f.POST("/comments", front.Comments)
	}

	// intermediate representation, handles
	// conversion of AST into an IR
	i := router.Group("/ir")
	{
		i.POST("/build", ir.Build)
	}

	// 'middle' of compiler stages, takes the
	// krug IR and type checks everything.
	m := router.Group("/mid")
	{
		b := m.Group("/build")
		{
			// module -> [build/scope] -> scope map.
			b.POST("/scope", middle.BuildScope)

			// module -> [build/scope_map] -> scope dict
			b.POST("/scope_dict", middle.BuildScopeDict)

			// module -> [build/type] -> type map.
			b.POST("/type", middle.BuildType)
		}

		m.POST("/unused_func", middle.UnusedFunctions)

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

	// backend of the compiler handles taking the
	// IR and generating code from it.
	b := router.Group("/back")
	{
		b.POST("/gen", back.Gen)
	}
}
