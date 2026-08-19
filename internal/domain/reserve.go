package domain

import "time"

type ReserveStatus string

const (
	ReserveStatusPending   ReserveStatus = "pending"
	ReserveStatusConfirmed ReserveStatus = "confirmed"
	ReserveStatusCanceled  ReserveStatus = "canceled"
	ReserveStatusExpired   ReserveStatus = "expired"
)

type Reserve struct {
	ID        int           `json:"id" gorm:"column:id;primaryKey"`
	ProductID int           `json:"product_id" gorm:"column:product_id"`
	Status    ReserveStatus `json:"status" gorm:"column:status"`
	Quantity  int           `json:"quantity" gorm:"column:quantity"`
	CreatedAt time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"column:updated_at"`
}

func (Reserve) TableName() string {
	return "reservations"
}
