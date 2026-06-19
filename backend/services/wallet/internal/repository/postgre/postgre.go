package postgre

import (
	"context"
	"fmt"
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

// Insert inserts a wallet to the database.
// If same data exists (user_id, currency), it will just return the record without insert or update.
func (w *Wallet) Insert(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
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

	if err == sdkpostgre.ErrNotFound {
		res, err = w.getWalletByUserIDAndCurrency(ctx, wallet.UserID, wallet.Currency)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-Insert] fail insert wallet with tx", "error", err)
		return nil, entity.ErrInternal
	}
	return convertDBWalletToEntityWallet(res), nil
}

func (w *Wallet) getWalletByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*db.Wallet, error) {
	res, err := w.queries.GetWalletByUserIdAndCurrency(ctx, db.GetWalletByUserIdAndCurrencyParams{
		UserID:   userID,
		Currency: currency,
	})

	if err == sdkpostgre.ErrNotFound {
		fmt.Println("here")
		return nil, entity.ErrEmptyWallet
	}
	if err != nil {
		slog.ErrorContext(ctx, "[PostgreWallet-getWalletByUserIDAndCurrency] internal error", "error", err)
		return nil, entity.ErrInternal
	}
	return res, nil
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
