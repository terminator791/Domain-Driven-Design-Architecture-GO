package promotion

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// Repository defines the interface for promotion persistence (Repository Pattern)
type Repository interface {
	Save(ctx context.Context, promotion *Promotion) error
	FindByID(ctx context.Context, id PromotionID) (*Promotion, error)
	FindActivePromotions(ctx context.Context) ([]*Promotion, error)
	FindPromotionsByType(ctx context.Context, promoType PromotionType) ([]*Promotion, error)
	FindExpiredPromotions(ctx context.Context) ([]*Promotion, error)
	Delete(ctx context.Context, id PromotionID) error
}

// CouponRepository defines the interface for coupon persistence
type CouponRepository interface {
	Save(ctx context.Context, coupon *Coupon) error
	FindByCode(ctx context.Context, code string) (*Coupon, error)
	FindByID(ctx context.Context, id CouponID) (*Coupon, error)
	FindByCustomerID(ctx context.Context, customerID order.CustomerID) ([]*Coupon, error)
	FindExpiredCoupons(ctx context.Context) ([]*Coupon, error)
	Delete(ctx context.Context, id CouponID) error
}

// EventPublisher defines the interface for publishing promotion domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}