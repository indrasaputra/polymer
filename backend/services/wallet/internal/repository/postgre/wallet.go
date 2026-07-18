package postgre

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/db"
	sdkpostgre "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
)

// Wallet is responsible to connect wallet entity with wallets table in PostgreSQL.
type Wallet struct {
	queries *db.Queries
}

// NewWallet creates an instance of Wallet.
func NewWallet(q *db.Queries) *Wallet {
	return &Wallet{queries: q}
}

// InsertWallet inserts a wallet to the database.
// If same data exists (user_id, currency), it will just return the record without insert or update.
func (w *Wallet) InsertWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	if wallet == nil {
		return nil, entity.ErrEmptyWallet
	}

	param := db.InsertWalletParams{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
		CreatedBy: wallet.UserID,
		UpdatedBy: wallet.UserID,
	}
	res, err := w.queries.InsertWallet(ctx, param)

	// if wallet exist, postgre returns not found, so we handle this by querying the wallet
	if err == sdkpostgre.ErrNotFound {
		res, err = w.getUserActiveWalletByUserIDAndCurrency(ctx, wallet.UserID, wallet.Currency)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-Insert] fail insert wallet", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBWalletToEntityWallet(res), nil
}

func (w *Wallet) getUserActiveWalletByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*db.Wallet, error) {
	res, err := w.queries.GetUserActiveWalletByUserIdAndCurrency(ctx, db.GetUserActiveWalletByUserIdAndCurrencyParams{
		UserID:   userID,
		Currency: currency,
	})

	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrEmptyWallet
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-getWalletByUserIDAndCurrency] internal error", "error", err)
		return nil, entity.ErrInternal
	}
	return res, nil
}

// GetCustomerByUserID gets a customer by user id.
func (w *Wallet) GetCustomerByUserID(ctx context.Context, userID uuid.UUID) (*entity.Customer, error) {
	res, err := w.queries.GetCustomerByUserID(ctx, userID)
	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrNilCustomer
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetCustomerByUserID] fail get customer", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBCustomerToEntityCustomer(res), nil
}

// InsertCustomer inserts a customer to the database.
// If same data exists (by user_id), it will just return the record without insert or update.
func (w *Wallet) InsertCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	if customer == nil {
		return nil, entity.ErrNilCustomer
	}

	param := db.InsertCustomerParams{
		ID:               customer.ID,
		UserID:           customer.UserID,
		StripeCustomerID: customer.StripeCustomerID,
		CreatedAt:        customer.CreatedAt,
		UpdatedAt:        customer.UpdatedAt,
		CreatedBy:        customer.UserID,
		UpdatedBy:        customer.UserID,
	}
	res, err := w.queries.InsertCustomer(ctx, param)

	// if customer exist, postgre returns not found, so we handle this by querying the customer
	if err == sdkpostgre.ErrNotFound {
		var c *entity.Customer
		c, err = w.GetCustomerByUserID(ctx, customer.UserID)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-InsertCustomer] fail insert customer", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBCustomerToEntityCustomer(res), nil
}

func (w *Wallet) GetPendingTransactionByIdempotencyKey(ctx context.Context, key uuid.UUID) (*entity.Transaction, error) {
	res, err := w.queries.GetPendingTransactionByIdempotencyKey(ctx, key)
	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrNilTransaction
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetPendingTransactionByIdempotencyKey] fail get transaction", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBTransactionToEntityTransaction(res), nil
}

func convertDBWalletToEntityWallet(w *db.Wallet) *entity.Wallet {
	return &entity.Wallet{
		ID:       w.ID,
		UserID:   w.UserID,
		Currency: w.Currency,
		Balance:  w.Balance,
		Auditable: entity.Auditable{
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
			DeletedAt: w.DeletedAt,
			CreatedBy: w.CreatedBy,
			UpdatedBy: w.UpdatedBy,
			DeletedBy: w.DeletedBy,
		},
	}
}

func convertDBCustomerToEntityCustomer(c *db.Customer) *entity.Customer {
	return &entity.Customer{
		ID:               c.ID,
		UserID:           c.UserID,
		StripeCustomerID: c.StripeCustomerID,
		Auditable: entity.Auditable{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
			CreatedBy: c.CreatedBy,
			UpdatedBy: c.UpdatedBy,
			DeletedBy: c.DeletedBy,
		},
	}
}

func convertDBTransactionToEntityTransaction(t *db.Transaction) *entity.Transaction {
	return &entity.Transaction{
		ID:               t.ID,
		UserID:           t.UserID,
		Type:             entity.TransactionType(t.Type),
		Status:           entity.TransactionStatus(t.Status),
		IdempotencyKey:   t.IdempotencyKey,
		Amount:           t.Amount,
		Currency:         t.Currency,
		PaymentSessionID: t.PaymentSessionID,
		Auditable: entity.Auditable{
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
			DeletedAt: t.DeletedAt,
			CreatedBy: t.CreatedBy,
			UpdatedBy: t.UpdatedBy,
			DeletedBy: t.DeletedBy,
		},
	}
}
