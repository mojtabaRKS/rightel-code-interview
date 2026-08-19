package order

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
	"github.com/mojtabaRKS/rightel-code-interview/internal/repository"
	svc "github.com/mojtabaRKS/rightel-code-interview/internal/service/order"
)

type OrderHandler struct {
	orderService orderService
}

func NewOrderHandler(orderService orderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type orderService interface {
	Reserve(ctx context.Context, productID, qty int) (*domain.Reserve, error)
	GetReservation(ctx context.Context, reservationID int) (*domain.Reserve, error)
	Confirm(ctx context.Context, reservationID int) (*domain.Reserve, error)
	Cancel(ctx context.Context, reservationID int) (*domain.Reserve, error)
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


func reservationIDFromPath(c *gin.Context) (int, bool) {
	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil || reservationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reservation id must be a positive integer"})
		return 0, false
	}

	return reservationID, true
}

