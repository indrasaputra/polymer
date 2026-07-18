package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/shopspring/decimal"
)

// TopupWallet defines interface to topup wallet.
type TopupWallet interface {
	// Topup topups a wallet's balance.
	// It needs idempotency key.
	Topup(ctx context.Context, input *entity.TopupWalletInput) (*entity.TopupWalletOutput, error)
}

// TopupWalletRepository defines the interface to update wallet in repository.
type TopupWalletRepository interface {
}

// WalletTopup is responsible for topup a wallet.
type WalletTopup struct {
	walletRepo TopupWalletRepository
}

// NewWalletTopup creates an instance of WalletTopup.
func NewWalletTopup(w TopupWalletRepository) *WalletTopup {
	return &WalletTopup{walletRepo: w}
}

func (wt *WalletTopup) Topup(ctx context.Context, input *entity.TopupWalletInput) (*entity.TopupWalletOutput, error) {
	if err := validateTopupWalletInput(input); err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] topup input is invalid", "error", err)
		return nil, err
	}
	return nil, nil
}

func validateTopupWalletInput(input *entity.TopupWalletInput) error {
	if input == nil {
		return entity.ErrEmptyInput
	}
	if input.UserID == uuid.Nil {
		return entity.ErrInvalidUser
	}
	if input.WalletID == uuid.Nil {
		return entity.ErrInvalidWallet
	}
	if input.IdempotencyKey == uuid.Nil {
		return entity.ErrInvalidIdempotencyKey
	}
	if input.Amount.LessThanOrEqual(decimal.Zero) {
		return entity.ErrInvalidTopupAmount
	}
	return nil
}
