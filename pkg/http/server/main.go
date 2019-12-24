package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yishanzhilu/api-template/pkg/http/server/middleware"
)

// Start will start a gin server
func Start() {
	server := gin.New()

	server.Use(middleware.RequestID())
	server.Use(middleware.GinLogger())
	server.Use(gin.Recovery())

	healthCheck(server)
	server.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	server.GET("/fatal", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	server.Use(middleware.Authenticate())
	server.GET("/secret", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	server.Run(viper.GetString("ADDR"))
}
