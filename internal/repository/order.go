package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const reservationTTL = 15 * time.Minute

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Reserve(ctx context.Context, productID, qty int) (*domain.Reserve, error) {
	if qty <= 0 {
		return nil, ErrInvalidQuantity
	}

	var reservation domain.Reserve
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
			UPDATE products
			SET stock = stock - ?, updated_at = NOW()
			WHERE id = ? AND stock >= ?
		`, qty, productID, qty)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			var count int64
			if err := tx.Model(&domain.Product{}).Where("id = ?", productID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return ErrProductNotFound
			}
			return ErrInsufficientStock
		}

		reservation = domain.Reserve{
			ProductID: productID,
			Status:    domain.ReserveStatusPending,
			Quantity:  qty,
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *orderRepository) GetReservation(ctx context.Context, reservationID int) (*domain.Reserve, error) {
	var reservation domain.Reserve
	if err := r.db.WithContext(ctx).First(&reservation, reservationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	return &reservation, nil
}

func (r *orderRepository) Confirm(ctx context.Context, reservationID int) (*domain.Reserve, error) {
	var reservation domain.Reserve
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&reservation, reservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrReservationNotFound
			}
			return err
		}

		if reservation.Status != domain.ReserveStatusPending {
			return ErrInvalidReservationTransition
		}

		var databaseTime time.Time
		if err := tx.Raw("SELECT NOW()").Scan(&databaseTime).Error; err != nil {
			return err
		}
		if !reservation.CreatedAt.Add(reservationTTL).After(databaseTime) {
			return ErrReservationExpired
		}

		result := tx.Model(&domain.Reserve{}).
			Where("id = ? AND status = ?", reservation.ID, domain.ReserveStatusPending).
			Updates(map[string]any{
				"status":     domain.ReserveStatusConfirmed,
				"updated_at": gorm.Expr("NOW()"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrInvalidReservationTransition
		}

		reservation.Status = domain.ReserveStatusConfirmed
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *orderRepository) Cancel(ctx context.Context, reservationID int) (*domain.Reserve, error) {
	var reservation domain.Reserve
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&reservation, reservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrReservationNotFound
			}
			return err
		}

		if reservation.Status != domain.ReserveStatusPending {
			return ErrInvalidReservationTransition
		}

		result := tx.Model(&domain.Reserve{}).
			Where("id = ? AND status = ?", reservation.ID, domain.ReserveStatusPending).
			Updates(map[string]any{
				"status":     domain.ReserveStatusCanceled,
				"updated_at": gorm.Expr("NOW()"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrInvalidReservationTransition
		}

		result = tx.Exec(`
			UPDATE products
			SET stock = stock + ?, updated_at = NOW()
			WHERE id = ?
		`, reservation.Quantity, reservation.ProductID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrProductNotFound
		}

		reservation.Status = domain.ReserveStatusCanceled
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *orderRepository) ExpirePending(ctx context.Context, limit int) error {
	if limit <= 0 {
		return ErrInvalidLimit
	}

	var reservationIDs []int
	if err := r.db.WithContext(ctx).
		Model(&domain.Reserve{}).
		Where("status = ? AND created_at + INTERVAL '15 minutes' <= NOW()", domain.ReserveStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Pluck("id", &reservationIDs).Error; err != nil {
		return err
	}

	for _, reservationID := range reservationIDs {
		if err := r.expireReservation(ctx, reservationID); err != nil {
			return err
		}
	}

	return nil
}

func (r *orderRepository) expireReservation(ctx context.Context, reservationID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reservation domain.Reserve
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&reservation, reservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		if reservation.Status != domain.ReserveStatusPending {
			return nil
		}

		var databaseTime time.Time
		if err := tx.Raw("SELECT NOW()").Scan(&databaseTime).Error; err != nil {
			return err
		}
		if reservation.CreatedAt.Add(reservationTTL).After(databaseTime) {
			return nil
		}

		result := tx.Model(&domain.Reserve{}).
			Where("id = ? AND status = ?", reservation.ID, domain.ReserveStatusPending).
			Updates(map[string]any{
				"status":     domain.ReserveStatusExpired,
				"updated_at": gorm.Expr("NOW()"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrInvalidReservationTransition
		}

		result = tx.Exec(`
			UPDATE products
			SET stock = stock + ?, updated_at = NOW()
			WHERE id = ?
		`, reservation.Quantity, reservation.ProductID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrProductNotFound
		}

		return nil
	})
}
