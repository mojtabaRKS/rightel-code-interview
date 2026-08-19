package domain

import "time"

type Product struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Stock     int       `gorm:"column:stock"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Product) TableName() string {
	return "products"
}
