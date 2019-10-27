package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	const VERSION = "0.0.1"
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"version": VERSION,
		})
	})

	services := []endpointService {
		&frontendService{},
		&middleService{},
		&backendService{},
		&irService{},
	}
	for _, service := range services {
		service.register(router)
	}
}
