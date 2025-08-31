package product

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// Product represents a product entity
type Product struct {
	id          order.ProductID
	name        string
	description string
	price       order.Money
	stockLevel  int
	isAvailable bool
}

func NewProduct(name, description string, price order.Money, stockLevel int) *Product {
	return &Product{
		id:          order.NewProductID(),
		name:        name,
		description: description,
		price:       price,
		stockLevel:  stockLevel,
		isAvailable: true,
	}
}

func (p *Product) ID() order.ProductID {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Description() string {
	return p.description
}

func (p *Product) Price() order.Money {
	return p.price
}

func (p *Product) StockLevel() int {
	return p.stockLevel
}

func (p *Product) IsAvailable() bool {
	return p.isAvailable
}

func (p *Product) UpdatePrice(price order.Money) {
	p.price = price
}

func (p *Product) UpdateStock(level int) {
	p.stockLevel = level
	p.isAvailable = level > 0
}

func (p *Product) HasSufficientStock(quantity order.Quantity) bool {
	return p.isAvailable && p.stockLevel >= quantity.Value()
}

// Repository interface for product persistence
type Repository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id order.ProductID) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	Delete(ctx context.Context, id order.ProductID) error
}