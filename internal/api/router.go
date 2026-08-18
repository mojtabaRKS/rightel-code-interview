package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes
// @title						scaffold service
// @version         			1.0.0
// @description     			This APIs create server for scaffolding services
// @Host 						localhost:8080
// @BasePath  					/
// @Schemes 					https
func (s *Server) SetupAPIRoutes() {
	r := s.engine

	{
		r.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		})
	}
}
