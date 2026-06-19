package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

// CreateWallet defines interface to create wallet.
type CreateWallet interface {
	// Create creates a new wallet.
	Create(ctx context.Context, input *entity.CreateWalletInput) (*entity.Wallet, error)
}

// CreateWalletRepository defines the interface to insert wallet to repository.
type CreateWalletRepository interface {
	// Insert inserts a wallet.
	Insert(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
}

// WalletCreator is responsible for creating a new wallet.
type WalletCreator struct {
	walletRepo CreateWalletRepository
}

// NewWalletCreator creates an instance of WalletCreator.
func NewWalletCreator(r CreateWalletRepository) *WalletCreator {
	return &WalletCreator{walletRepo: r}
}

// Create creates a new wallet.
// It is idempotent. If user already has wallet with the same currency as input,
// it will not create new wallet.
func (wc *WalletCreator) Create(ctx context.Context, input *entity.CreateWalletInput) (*entity.Wallet, error) {
	if err := validateCreateWalletInput(input); err != nil {
		slog.ErrorContext(ctx, "[WalletCreator-Create] wallet is invalid", "error", err)
		return nil, err
	}

	wallet := convertCreateWalletInputToWallet(input)

	result, err := wc.walletRepo.Insert(ctx, wallet)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletCreator-Create] fail save to repository", "error", err)
		return nil, err
	}
	return result, nil
}

func validateCreateWalletInput(wallet *entity.CreateWalletInput) error {
	if wallet == nil {
		return entity.ErrEmptyWallet
	}
	if wallet.UserID == uuid.Nil {
		return entity.ErrInvalidUser
	}
	// TODO: validate currency against collection

	return nil
}

func convertCreateWalletInputToWallet(input *entity.CreateWalletInput) *entity.Wallet {
	wallet := &entity.Wallet{
		ID:       uuid.Must(uuid.NewV7()),
		UserID:   input.UserID,
		Balance:  decimal.Zero,
		Currency: input.Currency,
	}
	setWalletAuditableProperties(wallet)
	return wallet
}

func setWalletAuditableProperties(wallet *entity.Wallet) {
	wallet.CreatedAt = time.Now().UTC()
	wallet.UpdatedAt = time.Now().UTC()
	wallet.CreatedBy = wallet.UserID
	wallet.UpdatedBy = wallet.UserID
}
