package common

import "errors"

var (
	// Domain errors
	ErrInvalidInput             = errors.New("invalid input")
	ErrEntityNotFound           = errors.New("entity not found")
	ErrInsufficientStock        = errors.New("insufficient stock")
	ErrInvalidPrice             = errors.New("invalid price")
	ErrInvalidQuantity          = errors.New("invalid quantity")
	ErrOrderNotFound            = errors.New("order not found")
	ErrProductNotFound          = errors.New("product not found")
	ErrCustomerNotFound         = errors.New("customer not found")
	
	// Cart specific errors
	ErrCartNotFound             = errors.New("cart not found")
	ErrCartExpired              = errors.New("cart expired")
	ErrProductNotAvailable      = errors.New("product not available")
	ErrMaxQuantityExceeded      = errors.New("maximum quantity per product exceeded")
	ErrCartTotalItemsExceeded   = errors.New("maximum total items in cart exceeded")
	
	// Inventory specific errors
	ErrReservationNotFound      = errors.New("stock reservation not found")
	ErrReservationExpired       = errors.New("stock reservation expired")
	ErrStockNotAvailable        = errors.New("stock not available for reservation")
	
	// Payment specific errors
	ErrInvalidPaymentStatus     = errors.New("invalid payment status for operation")
	ErrInvalidAmount            = errors.New("invalid payment amount")
	ErrRefundAmountExceedsCapture = errors.New("refund amount exceeds captured amount")
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentMethodNotFound    = errors.New("payment method not found")
	
	// Promotion specific errors
	ErrPromotionNotFound        = errors.New("promotion not found")
	ErrPromotionNotValid        = errors.New("promotion is not valid or active")
	ErrPromotionAlreadyActive   = errors.New("promotion is already active")
	ErrInvalidPromotionStatus   = errors.New("invalid promotion status for operation")
	ErrInvalidDateRange         = errors.New("invalid date range")
	ErrCouponNotFound           = errors.New("coupon not found")
	ErrCouponAlreadyUsed        = errors.New("coupon has already been used")
	ErrCouponExpired            = errors.New("coupon has expired")
	ErrCouponCodeExists         = errors.New("coupon code already exists")
)