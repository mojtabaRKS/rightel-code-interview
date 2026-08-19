package product

import (
	"context"
	"errors"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
)

const MaxPageSize = 100

var ErrInvalidPagination = errors.New("invalid pagination")

type productRepository interface {
	Search(ctx context.Context, query string, page, size int) ([]domain.Product, error)
}

type productService struct {
	productRepository productRepository
}

func NewProductService(productRepository productRepository) *productService {
	return &productService{productRepository: productRepository}
}

