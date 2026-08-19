package order

import (
	"context"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
)

func (s *orderService) Reserve(ctx context.Context, productID, qty int) (*domain.Reserve, error) {
	if productID <= 0 || qty <= 0 {
		return nil, ErrInvalidInput
	}

	reservation, err := s.orderRepository.Reserve(ctx, productID, qty)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *orderService) GetReservation(ctx context.Context, reservationID int) (*domain.Reserve, error) {
	if reservationID <= 0 {
		return nil, ErrInvalidInput
	}

	return s.orderRepository.GetReservation(ctx, reservationID)
}

func (s *orderService) Confirm(ctx context.Context, reservationID int) (*domain.Reserve, error) {
	if reservationID <= 0 {
		return nil, ErrInvalidInput
	}

	return s.orderRepository.Confirm(ctx, reservationID)
}

func (s *orderService) Cancel(ctx context.Context, reservationID int) (*domain.Reserve, error) {
	if reservationID <= 0 {
		return nil, ErrInvalidInput
	}

	return s.orderRepository.Cancel(ctx, reservationID)
}

func (s *orderService) ExpirePending(ctx context.Context, limit int) error {
	if limit <= 0 {
		return ErrInvalidInput
	}

	return s.orderRepository.ExpirePending(ctx, limit)
}
