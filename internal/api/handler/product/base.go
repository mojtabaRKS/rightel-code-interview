package product

import (
	"context"

	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
)

type ProductHandler struct {
	productService productService
}

func NewProductHandler(productService productService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type productService interface {
	Search(ctx context.Context, query string, page, size int) ([]domain.Product, error)
}
