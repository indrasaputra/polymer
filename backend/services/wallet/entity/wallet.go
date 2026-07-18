package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ContextKey is just a typed-string.
type ContextKey string

// ContextKeyCurrentUser should be used as key in context.
const ContextKeyCurrentUser = "CURRENT_USER"

// CurrentUser represents current user from JWT.
type CurrentUser struct {
	Email string
	ID    uuid.UUID
}

// CreateWalletInput defines logical data for create wallet.
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

// TopupWalletInput defines logical input data related to topup wallet.
type TopupWalletInput struct {
	Amount         decimal.Decimal
	WalletID       uuid.UUID
	UserID         uuid.UUID
	IdempotencyKey uuid.UUID
}

// TopupWalletOutput defines logical output data related to topup wallet.
type TopupWalletOutput struct {
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
