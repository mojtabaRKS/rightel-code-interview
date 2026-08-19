package order

import (
	"errors"
	"net/http"
	"strconv"

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

func (oh *OrderHandler) GetReservation(c *gin.Context) {
	reservationID, ok := reservationIDFromPath(c)
	if !ok {
		return
	}

	reservation, err := oh.orderService.GetReservation(c.Request.Context(), reservationID)
	if err != nil {
		respondLifecycleError(c, err)
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func (oh *OrderHandler) Confirm(c *gin.Context) {
	reservationID, ok := reservationIDFromPath(c)
	if !ok {
		return
	}

	reservation, err := oh.orderService.Confirm(c.Request.Context(), reservationID)
	if err != nil {
		respondLifecycleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation_id": reservation.ID,
		"status":         reservation.Status,
	})
}

func (oh *OrderHandler) Cancel(c *gin.Context) {
	reservationID, ok := reservationIDFromPath(c)
	if !ok {
		return
	}

	reservation, err := oh.orderService.Cancel(c.Request.Context(), reservationID)
	if err != nil {
		respondLifecycleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation_id": reservation.ID,
		"status":         reservation.Status,
	})
}

func reservationIDFromPath(c *gin.Context) (int, bool) {
	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil || reservationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reservation id must be a positive integer"})
		return 0, false
	}

	return reservationID, true
}

func respondLifecycleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, svc.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
	case errors.Is(err, repository.ErrReservationNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
	case errors.Is(err, repository.ErrReservationExpired):
		c.JSON(http.StatusConflict, gin.H{"error": "reservation expired"})
	case errors.Is(err, repository.ErrInvalidReservationTransition):
		c.JSON(http.StatusConflict, gin.H{"error": "invalid reservation transition"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
