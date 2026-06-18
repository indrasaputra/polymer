package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateWalletInput defines logical data for create wallet.
type CreateWalletInput struct {
	Currency string    `json:"currency"`
	UserID   uuid.UUID `json:"user_id"`
}

// Wallet defines logical data related to wallet.
type Wallet struct {
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"currency"`
	Auditable
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
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
