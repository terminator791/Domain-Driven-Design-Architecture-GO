package common

import "errors"

var (
	// Domain errors
	ErrInvalidInput      = errors.New("invalid input")
	ErrEntityNotFound    = errors.New("entity not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidPrice      = errors.New("invalid price")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrOrderNotFound     = errors.New("order not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrCustomerNotFound  = errors.New("customer not found")
)