package promotion

import (
	"time"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// Coupon represents a discount coupon (Entity)
type Coupon struct {
	id           CouponID
	code         string
	promotionID  PromotionID
	customerID   *order.CustomerID // nil means available to all customers
	isUsed       bool
	usedAt       *time.Time
	expiresAt    time.Time
	createdAt    time.Time
}

func NewCoupon(code string, promotionID PromotionID, customerID *order.CustomerID, expiresAt time.Time) *Coupon {
	return &Coupon{
		id:          NewCouponID(),
		code:        code,
		promotionID: promotionID,
		customerID:  customerID,
		isUsed:      false,
		expiresAt:   expiresAt,
		createdAt:   time.Now(),
	}
}

func (c *Coupon) ID() CouponID {
	return c.id
}

func (c *Coupon) Code() string {
	return c.code
}

func (c *Coupon) PromotionID() PromotionID {
	return c.promotionID
}

func (c *Coupon) CustomerID() *order.CustomerID {
	return c.customerID
}

func (c *Coupon) IsUsed() bool {
	return c.isUsed
}

func (c *Coupon) UsedAt() *time.Time {
	return c.usedAt
}

func (c *Coupon) ExpiresAt() time.Time {
	return c.expiresAt
}

func (c *Coupon) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Coupon) IsExpired() bool {
	return time.Now().After(c.expiresAt)
}

func (c *Coupon) IsValidForCustomer(customerID order.CustomerID) bool {
	if c.customerID != nil && !c.customerID.Equals(customerID) {
		return false
	}
	return !c.isUsed && !c.IsExpired()
}

func (c *Coupon) Use() error {
	if c.isUsed {
		return common.ErrCouponAlreadyUsed
	}
	if c.IsExpired() {
		return common.ErrCouponExpired
	}
	
	c.isUsed = true
	now := time.Now()
	c.usedAt = &now
	return nil
}

// Promotion represents a promotional campaign (Aggregate Root Entity)
type Promotion struct {
	id          PromotionID
	name        string
	description string
	promoType   PromotionType
	discount    DiscountValue
	conditions  PromotionCondition
	usageLimits UsageLimits
	status      PromotionStatus
	startsAt    time.Time
	endsAt      time.Time
	createdAt   time.Time
	updatedAt   time.Time
	createdBy   string
	coupons     []Coupon
	usageHistory []PromotionUsage
	events      []interface{} // Domain events
}

func NewPromotion(name, description string, promoType PromotionType, discount DiscountValue, createdBy string) *Promotion {
	now := time.Now()
	return &Promotion{
		id:           NewPromotionID(),
		name:         name,
		description:  description,
		promoType:    promoType,
		discount:     discount,
		conditions:   NewPromotionCondition(),
		usageLimits:  NewUsageLimits(0, 0), // Unlimited by default
		status:       PromotionStatusDraft,
		createdAt:    now,
		updatedAt:    now,
		createdBy:    createdBy,
		coupons:      make([]Coupon, 0),
		usageHistory: make([]PromotionUsage, 0),
		events:       make([]interface{}, 0),
	}
}

func (p *Promotion) ID() PromotionID {
	return p.id
}

func (p *Promotion) Name() string {
	return p.name
}

func (p *Promotion) Description() string {
	return p.description
}

func (p *Promotion) Type() PromotionType {
	return p.promoType
}

func (p *Promotion) Discount() DiscountValue {
	return p.discount
}

func (p *Promotion) Conditions() PromotionCondition {
	return p.conditions
}

func (p *Promotion) UsageLimits() UsageLimits {
	return p.usageLimits
}

func (p *Promotion) Status() PromotionStatus {
	return p.status
}

func (p *Promotion) StartsAt() time.Time {
	return p.startsAt
}

func (p *Promotion) EndsAt() time.Time {
	return p.endsAt
}

func (p *Promotion) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Promotion) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Promotion) CreatedBy() string {
	return p.createdBy
}

func (p *Promotion) Coupons() []Coupon {
	return p.coupons
}

func (p *Promotion) UsageHistory() []PromotionUsage {
	return p.usageHistory
}

func (p *Promotion) Events() []interface{} {
	return p.events
}

func (p *Promotion) ClearEvents() {
	p.events = make([]interface{}, 0)
}

// SetConditions sets the promotion conditions
func (p *Promotion) SetConditions(conditions PromotionCondition) {
	p.conditions = conditions
	p.updatedAt = time.Now()
}

// SetUsageLimits sets the usage limits
func (p *Promotion) SetUsageLimits(limits UsageLimits) {
	p.usageLimits = limits
	p.updatedAt = time.Now()
}

// SetActiveperiod sets the active period for the promotion
func (p *Promotion) SetActivePeriod(startsAt, endsAt time.Time) error {
	if startsAt.After(endsAt) {
		return common.ErrInvalidDateRange
	}
	
	p.startsAt = startsAt
	p.endsAt = endsAt
	p.updatedAt = time.Now()
	return nil
}

// Activate activates the promotion
func (p *Promotion) Activate() error {
	if p.status == PromotionStatusActive {
		return common.ErrPromotionAlreadyActive
	}
	
	if p.status == PromotionStatusExpired || p.status == PromotionStatusCancelled {
		return common.ErrInvalidPromotionStatus
	}
	
	p.status = PromotionStatusActive
	p.updatedAt = time.Now()
	return nil
}

// Pause pauses the promotion
func (p *Promotion) Pause() error {
	if p.status != PromotionStatusActive {
		return common.ErrInvalidPromotionStatus
	}
	
	p.status = PromotionStatusPaused
	p.updatedAt = time.Now()
	return nil
}

// Cancel cancels the promotion
func (p *Promotion) Cancel() error {
	if p.status == PromotionStatusExpired || p.status == PromotionStatusCancelled {
		return common.ErrInvalidPromotionStatus
	}
	
	p.status = PromotionStatusCancelled
	p.updatedAt = time.Now()
	return nil
}

// IsActive checks if the promotion is currently active
func (p *Promotion) IsActive() bool {
	now := time.Now()
	return p.status == PromotionStatusActive && 
		   (p.startsAt.IsZero() || now.After(p.startsAt)) && 
		   (p.endsAt.IsZero() || now.Before(p.endsAt))
}

// IsExpired checks if the promotion is expired
func (p *Promotion) IsExpired() bool {
	if !p.endsAt.IsZero() && time.Now().After(p.endsAt) {
		return true
	}
	return false
}

// CanUse checks if the promotion can be used based on usage limits
func (p *Promotion) CanUse() bool {
	return p.IsActive() && p.usageLimits.CanUse()
}

// GenerateCoupon generates a coupon for this promotion
func (p *Promotion) GenerateCoupon(code string, customerID *order.CustomerID, expiresAt time.Time) (*Coupon, error) {
	if p.status == PromotionStatusCancelled || p.status == PromotionStatusExpired {
		return nil, common.ErrInvalidPromotionStatus
	}
	
	// Check if coupon code already exists
	for _, coupon := range p.coupons {
		if coupon.code == code {
			return nil, common.ErrCouponCodeExists
		}
	}
	
	coupon := NewCoupon(code, p.id, customerID, expiresAt)
	p.coupons = append(p.coupons, *coupon)
	p.updatedAt = time.Now()
	
	return coupon, nil
}

// ApplyToOrder applies the promotion to an order
func (p *Promotion) ApplyToOrder(orderTotal order.Money, customerID order.CustomerID, orderID order.OrderID) (order.Money, error) {
	if !p.CanUse() {
		return order.Money{}, common.ErrPromotionNotValid
	}
	
	// Apply the discount
	discountAmount, err := p.discount.CalculateDiscount(orderTotal)
	if err != nil {
		return order.Money{}, err
	}
	
	// Record usage
	usage := PromotionUsage{
		customerID:     customerID,
		orderID:        orderID,
		discountAmount: discountAmount,
		usedAt:         time.Now(),
	}
	p.usageHistory = append(p.usageHistory, usage)
	p.usageLimits.IncrementUsage()
	p.updatedAt = time.Now()
	
	// Add domain event
	event := PromotionAppliedEvent{
		PromotionID:    p.id,
		CustomerID:     customerID,
		OrderID:        orderID,
		DiscountAmount: discountAmount,
		AppliedAt:      time.Now(),
	}
	p.events = append(p.events, event)
	
	return discountAmount, nil
}

// GetCustomerUsageCount returns how many times a customer has used this promotion
func (p *Promotion) GetCustomerUsageCount(customerID order.CustomerID) int {
	count := 0
	for _, usage := range p.usageHistory {
		if usage.customerID.Equals(customerID) {
			count++
		}
	}
	return count
}

// CanCustomerUse checks if a specific customer can use this promotion
func (p *Promotion) CanCustomerUse(customerID order.CustomerID) bool {
	if !p.CanUse() {
		return false
	}
	
	// Check per-customer usage limits
	if p.usageLimits.MaxUsesPerCustomer() > 0 {
		customerUsageCount := p.GetCustomerUsageCount(customerID)
		if customerUsageCount >= p.usageLimits.MaxUsesPerCustomer() {
			return false
		}
	}
	
	return true
}

// CheckExpiry checks if the promotion should be expired and updates status
func (p *Promotion) CheckExpiry() {
	if p.IsExpired() && p.status == PromotionStatusActive {
		p.status = PromotionStatusExpired
		p.updatedAt = time.Now()
		
		// Add domain event
		event := PromotionExpiredEvent{
			PromotionID: p.id,
			Name:        p.name,
			ExpiredAt:   time.Now(),
			TotalUses:   p.usageLimits.CurrentTotalUses(),
		}
		p.events = append(p.events, event)
	}
}

// PromotionUsage represents a record of promotion usage (Value Object)
type PromotionUsage struct {
	customerID     order.CustomerID
	orderID        order.OrderID
	discountAmount order.Money
	usedAt         time.Time
}

func (u PromotionUsage) CustomerID() order.CustomerID {
	return u.customerID
}

func (u PromotionUsage) OrderID() order.OrderID {
	return u.orderID
}

func (u PromotionUsage) DiscountAmount() order.Money {
	return u.discountAmount
}

func (u PromotionUsage) UsedAt() time.Time {
	return u.usedAt
}