package payment

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// PaymentID represents a unique identifier for a payment (Value Object)
type PaymentID struct {
	value string
}

func NewPaymentID() PaymentID {
	return PaymentID{value: uuid.New().String()}
}

func NewPaymentIDFromString(id string) (PaymentID, error) {
	if id == "" {
		return PaymentID{}, errors.New("payment ID cannot be empty")
	}
	return PaymentID{value: id}, nil
}

func (id PaymentID) String() string {
	return id.value
}

func (id PaymentID) Equals(other PaymentID) bool {
	return id.value == other.value
}

// PaymentMethodID represents a unique identifier for a payment method (Value Object)
type PaymentMethodID struct {
	value string
}

func NewPaymentMethodID() PaymentMethodID {
	return PaymentMethodID{value: uuid.New().String()}
}

func NewPaymentMethodIDFromString(id string) (PaymentMethodID, error) {
	if id == "" {
		return PaymentMethodID{}, errors.New("payment method ID cannot be empty")
	}
	return PaymentMethodID{value: id}, nil
}

func (id PaymentMethodID) String() string {
	return id.value
}

func (id PaymentMethodID) Equals(other PaymentMethodID) bool {
	return id.value == other.value
}

// PaymentStatus represents the status of a payment (Value Object)
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusAuthorized PaymentStatus = "AUTHORIZED"
	PaymentStatusCaptured  PaymentStatus = "CAPTURED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

func (s PaymentStatus) String() string {
	return string(s)
}

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusPending, PaymentStatusAuthorized, PaymentStatusCaptured, PaymentStatusRefunded, PaymentStatusFailed, PaymentStatusCancelled:
		return true
	default:
		return false
	}
}

// PaymentMethodType represents different types of payment methods (Value Object)
type PaymentMethodType string

const (
	PaymentMethodTypeCreditCard PaymentMethodType = "CREDIT_CARD"
	PaymentMethodTypeDebitCard  PaymentMethodType = "DEBIT_CARD"
	PaymentMethodTypePayPal     PaymentMethodType = "PAYPAL"
	PaymentMethodTypeBankTransfer PaymentMethodType = "BANK_TRANSFER"
	PaymentMethodTypeDigitalWallet PaymentMethodType = "DIGITAL_WALLET"
	PaymentMethodTypeCrypto     PaymentMethodType = "CRYPTOCURRENCY"
)

func (t PaymentMethodType) String() string {
	return string(t)
}

func (t PaymentMethodType) IsValid() bool {
	switch t {
	case PaymentMethodTypeCreditCard, PaymentMethodTypeDebitCard, PaymentMethodTypePayPal, PaymentMethodTypeBankTransfer, PaymentMethodTypeDigitalWallet, PaymentMethodTypeCrypto:
		return true
	default:
		return false
	}
}

// RefundReason represents the reason for a refund (Value Object)
type RefundReason string

const (
	RefundReasonCustomerRequest RefundReason = "CUSTOMER_REQUEST"
	RefundReasonDefectiveProduct RefundReason = "DEFECTIVE_PRODUCT"
	RefundReasonWrongItem       RefundReason = "WRONG_ITEM"
	RefundReasonNotAsDescribed  RefundReason = "NOT_AS_DESCRIBED"
	RefundReasonDamaged         RefundReason = "DAMAGED"
	RefundReasonSystemError     RefundReason = "SYSTEM_ERROR"
)

func (r RefundReason) String() string {
	return string(r)
}

// PaymentCreatedEvent represents a domain event when a payment is created
type PaymentCreatedEvent struct {
	PaymentID     PaymentID
	OrderID       order.OrderID
	CustomerID    order.CustomerID
	Amount        order.Money
	PaymentMethod PaymentMethodType
	CreatedAt     time.Time
}

// PaymentAuthorizedEvent represents a domain event when a payment is authorized
type PaymentAuthorizedEvent struct {
	PaymentID     PaymentID
	OrderID       order.OrderID
	Amount        order.Money
	AuthorizedAt  time.Time
	AuthCode      string
}

// PaymentCapturedEvent represents a domain event when a payment is captured
type PaymentCapturedEvent struct {
	PaymentID   PaymentID
	OrderID     order.OrderID
	Amount      order.Money
	CapturedAt  time.Time
	TransactionID string
}

// PaymentFailedEvent represents a domain event when a payment fails
type PaymentFailedEvent struct {
	PaymentID   PaymentID
	OrderID     order.OrderID
	Amount      order.Money
	FailedAt    time.Time
	Reason      string
	ErrorCode   string
}

// PaymentRefundedEvent represents a domain event when a payment is refunded
type PaymentRefundedEvent struct {
	PaymentID     PaymentID
	OrderID       order.OrderID
	RefundedAmount order.Money
	RefundedAt    time.Time
	Reason        RefundReason
	RefundID      string
}