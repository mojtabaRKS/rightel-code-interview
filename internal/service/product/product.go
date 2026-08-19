package product

import (
	"context"
	"math"
	"strings"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
)

func (s *productService) Search(ctx context.Context, query string, page, size int) ([]domain.Product, error) {
	if page <= 0 || size <= 0 || size > MaxPageSize || page > math.MaxInt/size {
		return nil, ErrInvalidPagination
	}

	return s.productRepository.Search(ctx, strings.TrimSpace(query), page, size)
}
