package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes
// @title						user Importer Service
// @version         			1.0.0
// @description     			This APIs create server for importing users and fetch them
// @Host 						localhost:8080
// @BasePath  					/
// @Schemes 					https
func (s *Server) SetupAPIRoutes() {
	r := s.engine

	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		})
	}
}
