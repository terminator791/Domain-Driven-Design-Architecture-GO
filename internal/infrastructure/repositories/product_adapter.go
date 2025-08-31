package repositories

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/product"
)

// ProductRepositoryAdapter adapts the product repository to the interface expected by the domain service
type ProductRepositoryAdapter struct {
	productRepo product.Repository
}

func NewProductRepositoryAdapter(productRepo product.Repository) *ProductRepositoryAdapter {
	return &ProductRepositoryAdapter{
		productRepo: productRepo,
	}
}

func (a *ProductRepositoryAdapter) FindByID(ctx context.Context, id order.ProductID) (order.ProductInfo, error) {
	product, err := a.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Product already implements the ProductInfo interface
	return product, nil
}