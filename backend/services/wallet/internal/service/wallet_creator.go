package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/bojanz/currency"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/uow"
)

// CreateWallet defines interface to create wallet.
type CreateWallet interface {
	// Create creates a new wallet.
	Create(ctx context.Context, input *entity.CreateWalletInput) (*entity.Wallet, error)
}

// CreateWalletRepository defines the interface to insert wallet to repository.
type CreateWalletRepository interface {
	// InsertWallet inserts a wallet.
	InsertWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
	// InsertCustomer inserts a customer.
	InsertCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	// GetCustomerByUserID gets customer. I decided to put it in wallet repository because the usage is closely
	// related with wallet case, not a separate flow.
	GetCustomerByUserID(ctx context.Context, userID uuid.UUID) (*entity.Customer, error)
}

// CreateCustomerClient defines the interface to create customer in 3rd party side.
type CreateCustomerClient interface {
	// CreateCustomer creates a new customer and returns ID from 3rd party.
	CreateCustomer(ctx context.Context, email string) (string, error)
}

// WalletCreator is responsible for creating a new wallet.
type WalletCreator struct {
	txManager      uow.TxManager
	walletRepo     CreateWalletRepository
	customerClient CreateCustomerClient
}

// NewWalletCreator creates an instance of WalletCreator.
func NewWalletCreator(m uow.TxManager, r CreateWalletRepository, c CreateCustomerClient) *WalletCreator {
	return &WalletCreator{txManager: m, walletRepo: r, customerClient: c}
}

// Create creates a new wallet.
// It is idempotent. If user already has wallet with the same currency as input,
// it will not create new wallet.
func (wc *WalletCreator) Create(ctx context.Context, input *entity.CreateWalletInput) (*entity.Wallet, error) {
	if err := validateCreateWalletInput(input); err != nil {
		slog.ErrorContext(ctx, "[WalletCreator-Create] wallet is invalid", "error", err)
		return nil, err
	}

	customer, err := wc.getOrCreateCustomer(ctx, input)
	if err != nil {
		return nil, err
	}

	wallet := convertCreateWalletInputToWallet(input)

	var result *entity.Wallet
	err = wc.txManager.Do(ctx, func(ctx context.Context) error {
		_, err = wc.walletRepo.InsertCustomer(ctx, customer)
		if err != nil {
			slog.ErrorContext(ctx, "[WalletCreator-Create] fail save customer to repository", "error", err)
			return err
		}

		result, err = wc.walletRepo.InsertWallet(ctx, wallet)
		if err != nil {
			slog.ErrorContext(ctx, "[WalletCreator-Create] fail save wallet to repository", "error", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, entity.ErrInternal
	}

	return result, nil
}

func (wc *WalletCreator) getOrCreateCustomer(ctx context.Context, input *entity.CreateWalletInput) (*entity.Customer, error) {
	customer, err := wc.walletRepo.GetCustomerByUserID(ctx, input.UserID)
	if err != nil && err != entity.ErrNilCustomer {
		slog.ErrorContext(ctx, "[WalletCreator-getOrCreateCustomer] fail get customer", "error", err)
		return nil, entity.ErrInternal
	}
	if err == entity.ErrNilCustomer {
		customerID, err := wc.customerClient.CreateCustomer(ctx, input.Email)
		if err != nil {
			slog.ErrorContext(ctx, "[WalletCreator-getOrCreateCustomer] fail create customer to client", "error", err)
			return nil, entity.ErrInternal
		}

		customer = &entity.Customer{
			ID:               uuid.Must(uuid.NewV7()),
			UserID:           input.UserID,
			StripeCustomerID: customerID,
		}
		setCustomerAuditableProperties(customer)
	}
	return customer, nil
}

func validateCreateWalletInput(wallet *entity.CreateWalletInput) error {
	if wallet == nil {
		return entity.ErrEmptyWallet
	}
	if wallet.UserID == uuid.Nil {
		return entity.ErrInvalidUser
	}
	if !currency.IsValid(wallet.Currency) {
		return entity.ErrInvalidCurrency
	}

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
	now := time.Now().UTC()
	wallet.CreatedAt = now
	wallet.UpdatedAt = now
	wallet.CreatedBy = wallet.UserID
	wallet.UpdatedBy = wallet.UserID
}

func setCustomerAuditableProperties(customer *entity.Customer) {
	now := time.Now().UTC()
	customer.CreatedAt = now
	customer.UpdatedAt = now
	customer.CreatedBy = customer.UserID
	customer.UpdatedBy = customer.UserID
}
