package payment

import (
	"time"

	"github.com/google/uuid"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// PaymentMethod represents a customer's payment method (Entity)
type PaymentMethod struct {
	id           PaymentMethodID
	customerID   order.CustomerID
	methodType   PaymentMethodType
	displayName  string
	isDefault    bool
	isActive     bool
	metadata     map[string]string // Store card last 4 digits, expiry, etc.
	createdAt    time.Time
	lastUsedAt   *time.Time
}

func NewPaymentMethod(customerID order.CustomerID, methodType PaymentMethodType, displayName string, metadata map[string]string) *PaymentMethod {
	return &PaymentMethod{
		id:          NewPaymentMethodID(),
		customerID:  customerID,
		methodType:  methodType,
		displayName: displayName,
		isDefault:   false,
		isActive:    true,
		metadata:    metadata,
		createdAt:   time.Now(),
	}
}

func (pm *PaymentMethod) ID() PaymentMethodID {
	return pm.id
}

func (pm *PaymentMethod) CustomerID() order.CustomerID {
	return pm.customerID
}

func (pm *PaymentMethod) MethodType() PaymentMethodType {
	return pm.methodType
}

func (pm *PaymentMethod) DisplayName() string {
	return pm.displayName
}

func (pm *PaymentMethod) IsDefault() bool {
	return pm.isDefault
}

func (pm *PaymentMethod) IsActive() bool {
	return pm.isActive
}

func (pm *PaymentMethod) Metadata() map[string]string {
	return pm.metadata
}

func (pm *PaymentMethod) CreatedAt() time.Time {
	return pm.createdAt
}

func (pm *PaymentMethod) LastUsedAt() *time.Time {
	return pm.lastUsedAt
}

func (pm *PaymentMethod) SetAsDefault() {
	pm.isDefault = true
}

func (pm *PaymentMethod) RemoveDefault() {
	pm.isDefault = false
}

func (pm *PaymentMethod) Deactivate() {
	pm.isActive = false
}

func (pm *PaymentMethod) UpdateLastUsed() {
	now := time.Now()
	pm.lastUsedAt = &now
}

// Payment represents a payment transaction (Aggregate Root Entity)
type Payment struct {
	id                PaymentID
	orderID           order.OrderID
	customerID        order.CustomerID
	paymentMethodID   PaymentMethodID
	amount            order.Money
	authorizedAmount  order.Money
	capturedAmount    order.Money
	refundedAmount    order.Money
	status            PaymentStatus
	authCode          string
	transactionID     string
	processorResponse map[string]string
	failureReason     string
	createdAt         time.Time
	authorizedAt      *time.Time
	capturedAt        *time.Time
	failedAt          *time.Time
	refunds           []Refund
	events            []interface{} // Domain events
}

func NewPayment(orderID order.OrderID, customerID order.CustomerID, paymentMethodID PaymentMethodID, amount order.Money) *Payment {
	return &Payment{
		id:                NewPaymentID(),
		orderID:           orderID,
		customerID:        customerID,
		paymentMethodID:   paymentMethodID,
		amount:            amount,
		authorizedAmount:  order.Money{},
		capturedAmount:    order.Money{},
		refundedAmount:    order.Money{},
		status:            PaymentStatusPending,
		processorResponse: make(map[string]string),
		createdAt:         time.Now(),
		refunds:           make([]Refund, 0),
		events:            make([]interface{}, 0),
	}
}

func (p *Payment) ID() PaymentID {
	return p.id
}

func (p *Payment) OrderID() order.OrderID {
	return p.orderID
}

func (p *Payment) CustomerID() order.CustomerID {
	return p.customerID
}

func (p *Payment) PaymentMethodID() PaymentMethodID {
	return p.paymentMethodID
}

func (p *Payment) Amount() order.Money {
	return p.amount
}

func (p *Payment) AuthorizedAmount() order.Money {
	return p.authorizedAmount
}

func (p *Payment) CapturedAmount() order.Money {
	return p.capturedAmount
}

func (p *Payment) RefundedAmount() order.Money {
	return p.refundedAmount
}

func (p *Payment) Status() PaymentStatus {
	return p.status
}

func (p *Payment) AuthCode() string {
	return p.authCode
}

func (p *Payment) TransactionID() string {
	return p.transactionID
}

func (p *Payment) ProcessorResponse() map[string]string {
	return p.processorResponse
}

func (p *Payment) FailureReason() string {
	return p.failureReason
}

func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Payment) AuthorizedAt() *time.Time {
	return p.authorizedAt
}

func (p *Payment) CapturedAt() *time.Time {
	return p.capturedAt
}

func (p *Payment) FailedAt() *time.Time {
	return p.failedAt
}

func (p *Payment) Refunds() []Refund {
	return p.refunds
}

func (p *Payment) Events() []interface{} {
	return p.events
}

func (p *Payment) ClearEvents() {
	p.events = make([]interface{}, 0)
}

// Authorize authorizes the payment
func (p *Payment) Authorize(authCode string, authorizedAmount order.Money, processorResponse map[string]string) error {
	if p.status != PaymentStatusPending {
		return common.ErrInvalidPaymentStatus
	}

	if authorizedAmount.Amount() > p.amount.Amount() {
		return common.ErrInvalidAmount
	}

	p.status = PaymentStatusAuthorized
	p.authCode = authCode
	p.authorizedAmount = authorizedAmount
	p.processorResponse = processorResponse
	now := time.Now()
	p.authorizedAt = &now

	// Add domain event
	event := PaymentAuthorizedEvent{
		PaymentID:    p.id,
		OrderID:      p.orderID,
		Amount:       authorizedAmount,
		AuthorizedAt: now,
		AuthCode:     authCode,
	}
	p.events = append(p.events, event)

	return nil
}

// Capture captures the authorized payment
func (p *Payment) Capture(captureAmount order.Money, transactionID string, processorResponse map[string]string) error {
	if p.status != PaymentStatusAuthorized {
		return common.ErrInvalidPaymentStatus
	}

	if captureAmount.Amount() > p.authorizedAmount.Amount() {
		return common.ErrInvalidAmount
	}

	p.status = PaymentStatusCaptured
	p.transactionID = transactionID
	p.capturedAmount = captureAmount
	p.processorResponse = processorResponse
	now := time.Now()
	p.capturedAt = &now

	// Add domain event
	event := PaymentCapturedEvent{
		PaymentID:     p.id,
		OrderID:       p.orderID,
		Amount:        captureAmount,
		CapturedAt:    now,
		TransactionID: transactionID,
	}
	p.events = append(p.events, event)

	return nil
}

// Fail marks the payment as failed
func (p *Payment) Fail(reason string, errorCode string, processorResponse map[string]string) error {
	if p.status == PaymentStatusCaptured || p.status == PaymentStatusRefunded {
		return common.ErrInvalidPaymentStatus
	}

	p.status = PaymentStatusFailed
	p.failureReason = reason
	p.processorResponse = processorResponse
	now := time.Now()
	p.failedAt = &now

	// Add domain event
	event := PaymentFailedEvent{
		PaymentID: p.id,
		OrderID:   p.orderID,
		Amount:    p.amount,
		FailedAt:  now,
		Reason:    reason,
		ErrorCode: errorCode,
	}
	p.events = append(p.events, event)

	return nil
}

// Refund creates a refund for the payment
func (p *Payment) Refund(refundAmount order.Money, reason RefundReason, refundID string) error {
	if p.status != PaymentStatusCaptured {
		return common.ErrInvalidPaymentStatus
	}

	// Calculate total refunded amount including this refund
	totalRefunded := p.refundedAmount.Amount() + refundAmount.Amount()
	if totalRefunded > p.capturedAmount.Amount() {
		return common.ErrRefundAmountExceedsCapture
	}

	// Create refund
	refund := NewRefund(p.id, refundAmount, reason, refundID)
	p.refunds = append(p.refunds, *refund)

	// Update refunded amount
	newRefundedAmount, err := order.NewMoney(totalRefunded, p.refundedAmount.Currency())
	if err != nil {
		return err
	}
	p.refundedAmount = newRefundedAmount

	// If fully refunded, update status
	if totalRefunded >= p.capturedAmount.Amount() {
		p.status = PaymentStatusRefunded
	}

	// Add domain event
	event := PaymentRefundedEvent{
		PaymentID:      p.id,
		OrderID:        p.orderID,
		RefundedAmount: refundAmount,
		RefundedAt:     time.Now(),
		Reason:         reason,
		RefundID:       refundID,
	}
	p.events = append(p.events, event)

	return nil
}

// Cancel cancels the payment
func (p *Payment) Cancel() error {
	if p.status == PaymentStatusCaptured || p.status == PaymentStatusRefunded {
		return common.ErrInvalidPaymentStatus
	}

	p.status = PaymentStatusCancelled
	return nil
}

// GetAvailableRefundAmount returns the amount that can still be refunded
func (p *Payment) GetAvailableRefundAmount() order.Money {
	availableAmount := p.capturedAmount.Amount() - p.refundedAmount.Amount()
	availableMoney, _ := order.NewMoney(availableAmount, p.capturedAmount.Currency())
	return availableMoney
}

// IsFullyRefunded checks if the payment is fully refunded
func (p *Payment) IsFullyRefunded() bool {
	return p.refundedAmount.Amount() >= p.capturedAmount.Amount()
}

// Refund represents a refund transaction (Entity)
type Refund struct {
	id           string
	paymentID    PaymentID
	amount       order.Money
	reason       RefundReason
	refundID     string // External processor refund ID
	processedAt  time.Time
	processorResponse map[string]string
}

func NewRefund(paymentID PaymentID, amount order.Money, reason RefundReason, refundID string) *Refund {
	return &Refund{
		id:                uuid.New().String(),
		paymentID:         paymentID,
		amount:            amount,
		reason:            reason,
		refundID:          refundID,
		processedAt:       time.Now(),
		processorResponse: make(map[string]string),
	}
}

func (r *Refund) ID() string {
	return r.id
}

func (r *Refund) PaymentID() PaymentID {
	return r.paymentID
}

func (r *Refund) Amount() order.Money {
	return r.amount
}

func (r *Refund) Reason() RefundReason {
	return r.reason
}

func (r *Refund) RefundID() string {
	return r.refundID
}

func (r *Refund) ProcessedAt() time.Time {
	return r.processedAt
}

func (r *Refund) ProcessorResponse() map[string]string {
	return r.processorResponse
}