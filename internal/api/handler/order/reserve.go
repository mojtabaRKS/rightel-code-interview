package order

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojtabaRKS/rightel-code-interview/internal/repository"
	svc "github.com/mojtabaRKS/rightel-code-interview/internal/service/order"
)

func (oh *OrderHandler) Reserve(c *gin.Context) {
	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.ProductID <= 0 || req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id and quantity must be positive"})
		return
	}

	reservation, err := oh.orderService.Reserve(c.Request.Context(), req.ProductID, req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, svc.ErrInvalidInput), errors.Is(err, repository.ErrInvalidQuantity):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		case errors.Is(err, repository.ErrInsufficientStock):
			c.JSON(http.StatusConflict, gin.H{"error": "insufficient stock"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation_id": reservation.ID,
		"status":         reservation.Status,
	})
}