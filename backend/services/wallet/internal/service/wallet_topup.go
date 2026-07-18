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
	// GetPendingTransactionByIdempotencyKey gets a pending transaction by idempotency key.
	GetPendingTransactionByIdempotencyKey(ctx context.Context, key uuid.UUID) (*entity.Transaction, error)
}

// PaymentClient defines interface for payment.
type PaymentClient interface {
	// GetCheckoutSessionURL gets a checkout session url by id.
	GetCheckoutSessionURL(ctx context.Context, id string) (string, error)
}

// WalletTopup is responsible for topup a wallet.
type WalletTopup struct {
	walletRepo    TopupWalletRepository
	paymentClient PaymentClient
}

// NewWalletTopup creates an instance of WalletTopup.
func NewWalletTopup(w TopupWalletRepository, p PaymentClient) *WalletTopup {
	return &WalletTopup{walletRepo: w, paymentClient: p}
}

func (wt *WalletTopup) Topup(ctx context.Context, input *entity.TopupWalletInput) (*entity.TopupWalletOutput, error) {
	if err := validateTopupWalletInput(input); err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] topup input is invalid", "error", err)
		return nil, err
	}

	trx, err := wt.walletRepo.GetPendingTransactionByIdempotencyKey(ctx, input.IdempotencyKey)
	if err != nil && err != entity.ErrNilTransaction {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail get transaction", "error", err)
		return nil, entity.ErrInternal
	}
	// there is pending transaction with inputted idempotency key.
	// just return the payment checkout session.
	if trx != nil && trx.PaymentSessionID != nil {
		url, err := wt.paymentClient.GetCheckoutSessionURL(ctx, *trx.PaymentSessionID)
		if err != nil {
			slog.ErrorContext(ctx, "[WalletTopup-Topup] fail get checkout session", "error", err)
			return nil, entity.ErrInternal
		}
		return &entity.TopupWalletOutput{URL: url}, nil
	}

	// this flow below is for non-existent transaction

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
