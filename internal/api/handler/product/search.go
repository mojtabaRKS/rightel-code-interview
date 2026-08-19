package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mojtabaRKS/rightel-code-interview/internal/domain"
	productSvc "github.com/mojtabaRKS/rightel-code-interview/internal/service/product"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
)

type productResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

func (h *ProductHandler) Search(c *gin.Context) {
	page, err := positiveQueryInt(c, "page", defaultPage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive integer"})
		return
	}

	size, err := positiveQueryInt(c, "size", defaultPageSize)
	if err != nil || size > productSvc.MaxPageSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "size must be between 1 and 100"})
		return
	}

	products, err := h.productService.Search(c.Request.Context(), c.Query("search"), page, size)
	if err != nil {
		if errors.Is(err, productSvc.ErrInvalidPagination) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]productResponse, 0, len(products))
	for _, product := range products {
		response = append(response, newProductResponse(product))
	}

	c.JSON(http.StatusOK, response)
}

func positiveQueryInt(c *gin.Context, name string, defaultValue int) (int, error) {
	value := c.Query(name)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, productSvc.ErrInvalidPagination
	}

	return parsed, nil
}

func newProductResponse(product domain.Product) productResponse {
	return productResponse{
		ID:    product.ID,
		Name:  product.Name,
		Stock: product.Stock,
	}
}
