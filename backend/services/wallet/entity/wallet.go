package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type (
	// ContextKey is just a typed-string.
	ContextKey string
	// TransactionType defines transaction's type.
	TransactionType string
	// TransactionStatus defines transaction's status.
	TransactionStatus string
)

const (
	// ContextKeyCurrentUser should be used as key in context.
	ContextKeyCurrentUser = "CURRENT_USER"

	// TransactionTypeTopup is "topup".
	TransactionTypeTopup TransactionType = "topup"

	// TransactionStatusPending is "pending".
	TransactionStatusPending TransactionStatus = "pending"
	// TransactionStatusCompleted is "completed".
	TransactionStatusCompleted TransactionStatus = "completed"
	// TransactionStatusFailed is "pending".
	TransactionStatusFailed TransactionStatus = "failed"
	// TransactionStatusCancelled is "cancelled".
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// CurrentUser represents current user from JWT.
type CurrentUser struct {
	Email string
	ID    uuid.UUID
}

// CreateWalletInput defines logical input data for create wallet.
type CreateWalletInput struct {
	Currency string
	Email    string
	UserID   uuid.UUID
}

// Wallet defines logical data related to wallet.
type Wallet struct {
	Balance  decimal.Decimal
	Currency string
	Auditable
	ID     uuid.UUID
	UserID uuid.UUID
}

// Customer defines logical data related to customer.
type Customer struct {
	StripeCustomerID string
	Auditable
	ID     uuid.UUID
	UserID uuid.UUID
}

// Transaction defines logical data related to transaction.
type Transaction struct {
	CheckoutSessionID *string
	Type              TransactionType
	Status            TransactionStatus
	Amount            decimal.Decimal
	Currency          string
	Auditable
	ID             uuid.UUID
	UserID         uuid.UUID
	IdempotencyKey uuid.UUID
}

// TopupWalletInput defines logical input data related to topup wallet.
type TopupWalletInput struct {
	Amount         decimal.Decimal
	WalletID       uuid.UUID
	UserID         uuid.UUID
	IdempotencyKey uuid.UUID
}

// TopupWalletOutput defines logical output data related to topup wallet.
type TopupWalletOutput struct {
	CheckoutSessionURL string
}

// CheckoutInput defines logical input data related to checkout process.
type CheckoutInput struct {
	Amount           decimal.Decimal
	Currency         string
	StripeCustomerID string
	SuccessURL       string
	Purpose          string
	Quantity         int
	WalletID         uuid.UUID
}

// CheckoutSession defines logical data related to checkout session.
type CheckoutSession struct {
	ID  string
	URL string
}

// Auditable defines logical data related to audit.
type Auditable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}
