package order

import (
	"context"
	"errors"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
)

var ErrInvalidInput = errors.New("invalid input")

type orderService struct {
	orderRepository orderRepository
}

func NewOrderService(orderRepository orderRepository) *orderService {
	return &orderService{orderRepository: orderRepository}
}

type orderRepository interface {
	Reserve(ctx context.Context, productID, qty int) (*domain.Reserve, error)
	GetReservation(ctx context.Context, reservationID int) (*domain.Reserve, error)
	Confirm(ctx context.Context, reservationID int) (*domain.Reserve, error)
	Cancel(ctx context.Context, reservationID int) (*domain.Reserve, error)
	ExpirePending(ctx context.Context, limit int) error
}