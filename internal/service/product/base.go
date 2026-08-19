package product

import (
	"context"
	"errors"
	"math"
	"strings"

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

func (s *productService) Search(ctx context.Context, query string, page, size int) ([]domain.Product, error) {
	if page <= 0 || size <= 0 || size > MaxPageSize || page > math.MaxInt/size {
		return nil, ErrInvalidPagination
	}

	return s.productRepository.Search(ctx, strings.TrimSpace(query), page, size)
}
