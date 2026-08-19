package order

import (
	"context"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
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
