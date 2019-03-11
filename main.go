package main

import (
	"net/http"
	"os"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/krugc-api/back"
	"github.com/krug-lang/krugc-api/front"
	"github.com/krug-lang/krugc-api/ir"
	"github.com/krug-lang/krugc-api/middle"
)

func main() {
	raven.SetDSN(os.Getenv("SENTRY_KEY"))

	router := gin.Default()

	// compiler frontend, handles lexing/parsing
	f := router.Group("/front")
	{
		// lexical analysis
		f.POST("/lex", front.Tokenize)

		// parsing.
		f.POST("/parse", front.Parse)
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
		m.POST("/build_scope", middle.BuildScope)

		r := m.Group("/resolve")
		{
			r.POST("/type", middle.TypeResolve)
			r.POST("/symbol", middle.SymbolResolve)
		}

	}

	// backend of the compiler handles taking the
	// IR and generating code from it.
	b := router.Group("/back")
	{
		b.POST("/gen", back.Gen)
	}

	s := &http.Server{
		Addr:           "localhost:8001",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
