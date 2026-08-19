package repository

import (
	"context"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Search(ctx context.Context, query string, page, size int) ([]domain.Product, error) {
	products := make([]domain.Product, 0)
	db := r.db.WithContext(ctx).Model(&domain.Product{})
	if query != "" {
		db = db.Where("name ILIKE ?", "%"+query+"%")
	}

	if err := db.
		Order("id ASC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
