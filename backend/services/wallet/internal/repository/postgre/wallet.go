package postgre

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/db"
	sdkpostgre "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
)

var (
	systemUser, _ = uuid.Parse("00000000-0000-7575-8331-5757e3110000")
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
		return nil, entity.ErrWalletEmpty
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
		var r *entity.Wallet
		r, err = w.getActiveWalletByUserIDAndCurrency(ctx, wallet.UserID, wallet.Currency)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-Insert] fail insert wallet", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBWalletToEntityWallet(res), nil
}

func (w *Wallet) getActiveWalletByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*entity.Wallet, error) {
	res, err := w.queries.GetActiveWalletByUserIdAndCurrency(ctx, db.GetActiveWalletByUserIdAndCurrencyParams{
		UserID:   userID,
		Currency: currency,
	})

	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrWalletNotFound
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetActiveWalletByUserIdAndCurrency] internal error", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBWalletToEntityWallet(res), nil
}

// GetActiveWalletByIDAndUserID gets an active wallet.
func (w *Wallet) GetActiveWalletByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Wallet, error) {
	res, err := w.queries.GetActiveWalletByIDAndUserID(ctx, db.GetActiveWalletByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})

	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrWalletNotFound
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetActiveWalletByIDAndUserID] internal error", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBWalletToEntityWallet(res), nil
}

// GetActiveCustomerByUserID gets a customer by user id.
func (w *Wallet) GetActiveCustomerByUserID(ctx context.Context, userID uuid.UUID) (*entity.Customer, error) {
	res, err := w.queries.GetActiveCustomerByUserID(ctx, userID)
	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrCustomerNotFound
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetActiveCustomerByUserID] fail get customer", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBCustomerToEntityCustomer(res), nil
}

// InsertCustomer inserts a customer to the database.
// If same data exists (by user_id), it will just return the record without insert or update.
func (w *Wallet) InsertCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	if customer == nil {
		return nil, entity.ErrCustomerEmpty
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
		c, err = w.GetActiveCustomerByUserID(ctx, customer.UserID)
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

// InsertTransaction inserts a new transaction.
func (w *Wallet) InsertTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	if transaction == nil {
		return nil, entity.ErrTransactionEmpty
	}

	param := db.InsertTransactionParams{
		ID:                transaction.ID,
		UserID:            transaction.UserID,
		Type:              db.TransactionType(transaction.Type),
		IdempotencyKey:    transaction.IdempotencyKey,
		Status:            db.TransactionStatus(transaction.Status),
		Amount:            transaction.Amount,
		Currency:          transaction.Currency,
		CheckoutSessionID: transaction.CheckoutSessionID,
		CreatedAt:         transaction.CreatedAt,
		UpdatedAt:         transaction.UpdatedAt,
		CreatedBy:         transaction.UserID,
		UpdatedBy:         transaction.UserID,
	}

	trx, err := w.queries.InsertTransaction(ctx, param)
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-InsertTransaction] fail insert transaction", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBTransactionToEntityTransaction(trx), nil
}

// GetActivePendingTransactionByIdempotencyKey gets a pending transaction.
func (w *Wallet) GetActivePendingTransactionByIdempotencyKey(ctx context.Context, key uuid.UUID) (*entity.Transaction, error) {
	res, err := w.queries.GetActivePendingTransactionByIdempotencyKey(ctx, key)
	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrTransactionNotFound
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetActivePendingTransactionByIdempotencyKey] fail get transaction", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBTransactionToEntityTransaction(res), nil
}

// UpdateActiveTransactionToCompletedByCheckoutSessionID updates a pending transaction to completed.
func (w *Wallet) UpdateActiveTransactionToCompletedByCheckoutSessionID(ctx context.Context, id string) error {
	param := db.UpdateActiveTransactionToCompletedByCheckoutSessionIDParams{
		CheckoutSessionID: &id,
		UpdatedAt:         time.Now().UTC(),
		UpdatedBy:         systemUser,
	}

	_, err := w.queries.UpdateActiveTransactionToCompletedByCheckoutSessionID(ctx, param)
	if err == sdkpostgre.ErrNotFound {
		return entity.ErrTransactionNotFound
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-UpdateActiveTransactionToCompletedByCheckoutSessionID] fail update transaction", "error", err)
		return entity.ErrInternal
	}
	return nil
}

// GetActiveTransactionByCheckoutSessionIDForUpdate gets active transaction for update.
func (w *Wallet) GetActiveTransactionByCheckoutSessionIDForUpdate(ctx context.Context, sessionID string) (*entity.Transaction, error) {
	res, err := w.queries.GetActiveTransactionByCheckoutSessionIDForUpdate(ctx, &sessionID)
	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrTransactionNotFound
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetActiveTransactionByCheckoutSessionIDForUpdate] fail get transaction", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBTransactionToEntityTransaction(res), nil
}

// GetActiveWalletByIDForUpdate gets active wallet for update.
func (w *Wallet) GetActiveWalletByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Wallet, error) {
	res, err := w.queries.GetActiveWalletByIDForUpdate(ctx, id)
	if err == sdkpostgre.ErrNotFound {
		return nil, entity.ErrWalletNotFound
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-GetActiveWalletByIDForUpdate] fail get wallet", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBWalletToEntityWallet(res), nil
}

// AddActiveWalletBalance adds some amount to specific wallet.
func (w *Wallet) AddActiveWalletBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*entity.Wallet, error) {
	param := db.AddActiveWalletBalanceParams{ID: id, Amount: amount}
	res, err := w.queries.AddActiveWalletBalance(ctx, param)
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-AddActiveWalletBalance] internal error", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBWalletToEntityWallet(res), nil
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
		ID:                t.ID,
		UserID:            t.UserID,
		Type:              entity.TransactionType(t.Type),
		Status:            entity.TransactionStatus(t.Status),
		IdempotencyKey:    t.IdempotencyKey,
		Amount:            t.Amount,
		Currency:          t.Currency,
		CheckoutSessionID: t.CheckoutSessionID,
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
