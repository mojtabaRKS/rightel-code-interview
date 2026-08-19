package api

import (
	"github.com/mojtabaRKS/rightel-code-interview/internal/api/handler/order"
	"github.com/mojtabaRKS/rightel-code-interview/internal/api/handler/product"
)

// SetupAPIRoutes
// @title						scaffold service
// @version         			1.0.0
// @description     			This APIs create server for scaffolding services
// @Host 						localhost:8080
// @BasePath  					/
// @Schemes 					https
func (s *Server) SetupAPIRoutes(oh *order.OrderHandler, ph *product.ProductHandler) {
	r := s.engine

	{
		r.POST("/reservations", oh.Reserve)
		r.GET("/reservations/:id", oh.GetReservation)
		r.POST("/reservations/:id/confirm", oh.Confirm)
		r.POST("/reservations/:id/cancel", oh.Cancel)
		r.GET("/products", ph.Search)
	}
}
