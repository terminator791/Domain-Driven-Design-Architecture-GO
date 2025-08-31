package promotion

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// PromotionID represents a unique identifier for a promotion (Value Object)
type PromotionID struct {
	value string
}

func NewPromotionID() PromotionID {
	return PromotionID{value: uuid.New().String()}
}

func NewPromotionIDFromString(id string) (PromotionID, error) {
	if id == "" {
		return PromotionID{}, errors.New("promotion ID cannot be empty")
	}
	return PromotionID{value: id}, nil
}

func (id PromotionID) String() string {
	return id.value
}

func (id PromotionID) Equals(other PromotionID) bool {
	return id.value == other.value
}

// CouponID represents a unique identifier for a coupon (Value Object)
type CouponID struct {
	value string
}

func NewCouponID() CouponID {
	return CouponID{value: uuid.New().String()}
}

func NewCouponIDFromString(id string) (CouponID, error) {
	if id == "" {
		return CouponID{}, errors.New("coupon ID cannot be empty")
	}
	return CouponID{value: id}, nil
}

func (id CouponID) String() string {
	return id.value
}

func (id CouponID) Equals(other CouponID) bool {
	return id.value == other.value
}

// PromotionType represents different types of promotions (Value Object)
type PromotionType string

const (
	PromotionTypePercentageDiscount PromotionType = "PERCENTAGE_DISCOUNT"
	PromotionTypeFixedDiscount      PromotionType = "FIXED_DISCOUNT"
	PromotionTypeBuyXGetY           PromotionType = "BUY_X_GET_Y"
	PromotionTypeFreeShipping       PromotionType = "FREE_SHIPPING"
	PromotionTypeSpendXGetY         PromotionType = "SPEND_X_GET_Y"
	PromotionTypeBundleDiscount     PromotionType = "BUNDLE_DISCOUNT"
)

func (t PromotionType) String() string {
	return string(t)
}

func (t PromotionType) IsValid() bool {
	switch t {
	case PromotionTypePercentageDiscount, PromotionTypeFixedDiscount, PromotionTypeBuyXGetY, 
		 PromotionTypeFreeShipping, PromotionTypeSpendXGetY, PromotionTypeBundleDiscount:
		return true
	default:
		return false
	}
}

// DiscountValue represents a discount amount (Value Object)
type DiscountValue struct {
	amount   float64
	currency string
	isPercent bool
}

func NewPercentageDiscount(percentage float64) (DiscountValue, error) {
	if percentage < 0 || percentage > 100 {
		return DiscountValue{}, errors.New("percentage must be between 0 and 100")
	}
	return DiscountValue{
		amount:    percentage,
		isPercent: true,
	}, nil
}

func NewFixedDiscount(amount float64, currency string) (DiscountValue, error) {
	if amount < 0 {
		return DiscountValue{}, errors.New("discount amount cannot be negative")
	}
	if currency == "" {
		return DiscountValue{}, errors.New("currency cannot be empty")
	}
	return DiscountValue{
		amount:    amount,
		currency:  currency,
		isPercent: false,
	}, nil
}

func (d DiscountValue) Amount() float64 {
	return d.amount
}

func (d DiscountValue) Currency() string {
	return d.currency
}

func (d DiscountValue) IsPercentage() bool {
	return d.isPercent
}

func (d DiscountValue) CalculateDiscount(originalAmount order.Money) (order.Money, error) {
	if d.isPercent {
		discountAmount := originalAmount.Amount() * (d.amount / 100)
		return order.NewMoney(discountAmount, originalAmount.Currency())
	} else {
		if d.currency != originalAmount.Currency() {
			return order.Money{}, errors.New("currency mismatch for fixed discount")
		}
		discountAmount := d.amount
		if discountAmount > originalAmount.Amount() {
			discountAmount = originalAmount.Amount()
		}
		return order.NewMoney(discountAmount, d.currency)
	}
}

// PromotionStatus represents the status of a promotion (Value Object)
type PromotionStatus string

const (
	PromotionStatusDraft    PromotionStatus = "DRAFT"
	PromotionStatusActive   PromotionStatus = "ACTIVE"
	PromotionStatusPaused   PromotionStatus = "PAUSED"
	PromotionStatusExpired  PromotionStatus = "EXPIRED"
	PromotionStatusCancelled PromotionStatus = "CANCELLED"
)

func (s PromotionStatus) String() string {
	return string(s)
}

func (s PromotionStatus) IsValid() bool {
	switch s {
	case PromotionStatusDraft, PromotionStatusActive, PromotionStatusPaused, PromotionStatusExpired, PromotionStatusCancelled:
		return true
	default:
		return false
	}
}

// UsageLimits represents usage limitations for promotions (Value Object)
type UsageLimits struct {
	maxTotalUses     int // Maximum total uses across all customers
	maxUsesPerCustomer int // Maximum uses per customer
	currentTotalUses   int // Current total uses
}

func NewUsageLimits(maxTotal, maxPerCustomer int) UsageLimits {
	return UsageLimits{
		maxTotalUses:       maxTotal,
		maxUsesPerCustomer: maxPerCustomer,
		currentTotalUses:   0,
	}
}

func (u UsageLimits) MaxTotalUses() int {
	return u.maxTotalUses
}

func (u UsageLimits) MaxUsesPerCustomer() int {
	return u.maxUsesPerCustomer
}

func (u UsageLimits) CurrentTotalUses() int {
	return u.currentTotalUses
}

func (u UsageLimits) CanUse() bool {
	if u.maxTotalUses > 0 && u.currentTotalUses >= u.maxTotalUses {
		return false
	}
	return true
}

func (u *UsageLimits) IncrementUsage() {
	u.currentTotalUses++
}

// PromotionCondition represents conditions for promotion eligibility (Value Object)
type PromotionCondition struct {
	minOrderAmount    *order.Money
	eligibleProducts  []order.ProductID
	eligibleCategories []string
	excludedProducts  []order.ProductID
	requiredCustomerTier string
	firstTimeCustomersOnly bool
}

func NewPromotionCondition() PromotionCondition {
	return PromotionCondition{
		eligibleProducts:   make([]order.ProductID, 0),
		eligibleCategories: make([]string, 0),
		excludedProducts:   make([]order.ProductID, 0),
	}
}

func (c *PromotionCondition) SetMinOrderAmount(amount order.Money) {
	c.minOrderAmount = &amount
}

func (c *PromotionCondition) AddEligibleProduct(productID order.ProductID) {
	c.eligibleProducts = append(c.eligibleProducts, productID)
}

func (c *PromotionCondition) AddEligibleCategory(category string) {
	c.eligibleCategories = append(c.eligibleCategories, category)
}

func (c *PromotionCondition) AddExcludedProduct(productID order.ProductID) {
	c.excludedProducts = append(c.excludedProducts, productID)
}

func (c *PromotionCondition) SetCustomerTier(tier string) {
	c.requiredCustomerTier = tier
}

func (c *PromotionCondition) SetFirstTimeCustomersOnly(firstTimeOnly bool) {
	c.firstTimeCustomersOnly = firstTimeOnly
}

func (c PromotionCondition) MinOrderAmount() *order.Money {
	return c.minOrderAmount
}

func (c PromotionCondition) EligibleProducts() []order.ProductID {
	return c.eligibleProducts
}

func (c PromotionCondition) EligibleCategories() []string {
	return c.eligibleCategories
}

func (c PromotionCondition) ExcludedProducts() []order.ProductID {
	return c.excludedProducts
}

func (c PromotionCondition) RequiredCustomerTier() string {
	return c.requiredCustomerTier
}

func (c PromotionCondition) FirstTimeCustomersOnly() bool {
	return c.firstTimeCustomersOnly
}

// PromotionAppliedEvent represents a domain event when a promotion is applied
type PromotionAppliedEvent struct {
	PromotionID  PromotionID
	CustomerID   order.CustomerID
	OrderID      order.OrderID
	DiscountAmount order.Money
	CouponCode   string
	AppliedAt    time.Time
}

// CouponUsedEvent represents a domain event when a coupon is used
type CouponUsedEvent struct {
	CouponID     CouponID
	CouponCode   string
	CustomerID   order.CustomerID
	OrderID      order.OrderID
	PromotionID  PromotionID
	DiscountAmount order.Money
	UsedAt       time.Time
}

// PromotionCreatedEvent represents a domain event when a promotion is created
type PromotionCreatedEvent struct {
	PromotionID   PromotionID
	PromotionType PromotionType
	Name          string
	CreatedBy     string
	CreatedAt     time.Time
}

// PromotionExpiredEvent represents a domain event when a promotion expires
type PromotionExpiredEvent struct {
	PromotionID PromotionID
	Name        string
	ExpiredAt   time.Time
	TotalUses   int
}