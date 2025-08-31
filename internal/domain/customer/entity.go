package customer

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// Customer represents a customer entity
type Customer struct {
	id    order.CustomerID
	name  string
	email string
	phone string
}

func NewCustomer(name, email, phone string) *Customer {
	return &Customer{
		id:    order.NewCustomerID(),
		name:  name,
		email: email,
		phone: phone,
	}
}

func (c *Customer) ID() order.CustomerID {
	return c.id
}

func (c *Customer) Name() string {
	return c.name
}

func (c *Customer) Email() string {
	return c.email
}

func (c *Customer) Phone() string {
	return c.phone
}

func (c *Customer) UpdateContactInfo(email, phone string) {
	c.email = email
	c.phone = phone
}

// Repository interface for customer persistence
type Repository interface {
	Save(ctx context.Context, customer *Customer) error
	FindByID(ctx context.Context, id order.CustomerID) (*Customer, error)
	FindByEmail(ctx context.Context, email string) (*Customer, error)
	FindAll(ctx context.Context) ([]*Customer, error)
	Delete(ctx context.Context, id order.CustomerID) error
}