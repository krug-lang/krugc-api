package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/krug-lang/caasper/controller"
)

func main() {
	raven.SetDSN(os.Getenv("SENTRY_KEY"))

	router := gin.Default()
	controller.RegisterRoutes(router)

	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	ip := "localhost"
	if envIP := os.Getenv("IP"); envIP != "" {
		ip = envIP
	}

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", ip, port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Started on", s.Addr)
	s.ListenAndServe()
}
