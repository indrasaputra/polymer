package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v86"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/uow"
)

// HandleStripeEventRepository defines the interface for transaction-related table.
type HandleStripeEventRepository interface {
	// UpdateTransactionToCompletedByCheckoutSessionID updates the pending transaction to completed.
	UpdateTransactionToCompletedByCheckoutSessionID(ctx context.Context, sessionID string) error
	// GetActiveTransactionByCheckoutSessionIDForUpdate gets an active transaction by checkout session ID.
	GetActiveTransactionByCheckoutSessionIDForUpdate(ctx context.Context, sessionID string) (*entity.Transaction, error)
	// GetActiveUserWalletByIDForUpdate gets active user's wallet by ID.
	GetActiveWalletByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Wallet, error)
	// AddWalletBalance adds wallet's balance.
	AddWalletBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*entity.Wallet, error)
}

// StripeEventHandler is responsible to handle Stripe event.
type StripeEventHandler struct {
	txManager uow.TxManager
	repo      HandleStripeEventRepository
}

// NewStripeEventHandler creates an instance of StripeEventHandler.
func NewStripeEventHandler(t uow.TxManager, r HandleStripeEventRepository) *StripeEventHandler {
	return &StripeEventHandler{txManager: t, repo: r}
}

// HandleCheckoutSessionCompleted handles checkout.session.completed.
func (s *StripeEventHandler) HandleCheckoutSessionCompleted(ctx context.Context, session *stripe.CheckoutSession) error {
	if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		err := s.updatePendingTransactionToCompleted(ctx, session)
		if err != nil {
			slog.ErrorContext(ctx, "[StripeEventHandler-HandleCheckoutSessionCompleted] fail update pending transaction", "error", err)
			return err
		}
	}
	return nil
}

func (s *StripeEventHandler) getActiveWalletByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Wallet, error) {
	wallet, err := s.repo.GetActiveWalletByIDForUpdate(ctx, id)
	if err == entity.ErrNilWallet {
		slog.ErrorContext(ctx, "[StripeEventHandler-getActiveUserWalletForUpdate] wallet not found", "error", err)
		return nil, entity.ErrNilWallet
	}
	if err != nil {
		slog.ErrorContext(ctx, "[StripeEventHandler-getActiveUserWalletForUpdate] fail get wallet", "error", err)
		return nil, entity.ErrInternal
	}
	return wallet, nil
}

func (s *StripeEventHandler) getActiveTransactionByCheckoutSessionIDForUpdate(ctx context.Context, sessionID string) (*entity.Transaction, error) {
	trx, err := s.repo.GetActiveTransactionByCheckoutSessionIDForUpdate(ctx, sessionID)
	if err == entity.ErrNilTransaction {
		slog.ErrorContext(ctx, "[StripeEventHandler-getActiveTransactionByCheckoutSessionIDForUpdate] transaction not found", "error", err)
		return nil, entity.ErrNilTransaction
	}
	if err != nil {
		slog.ErrorContext(ctx, "[StripeEventHandler-getActiveTransactionByCheckoutSessionIDForUpdate] fail get transaction", "error", err)
		return nil, entity.ErrInternal
	}
	return trx, nil
}

// Why does wallet need to be locked here? Doesn't update wallet do the job correctly without locking?
// Update wallet locks the row. So, the step by step process is: lock transaction, update transaction, update wallet.
// The problem with this approach, there is implicit lock on wallet.
// The actual process is: lock transaction, update transaction, lock wallet, update wallet.
// Therefore, in any other process that includes both wallet and transaction,
// transaction MUST BE locked first, then lock wallet so it aligns with this function.
// If in other process, wallet is locked first then lock transaction, there is a chance for deadlock.
// In this function, the purpose of lock wallet explicitly is to create a mindset
// that in other process, wallet must be locked first before transaction.
// Most of the time, wallet is locked first (e.g: check if balance is sufficient) before locking the transaction.
func (s *StripeEventHandler) updatePendingTransactionToCompleted(ctx context.Context, session *stripe.CheckoutSession) error {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		id, err := uuid.Parse(session.ClientReferenceID)
		if err != nil {
			return entity.ErrBadRequest
		}

		wallet, err := s.getActiveWalletByIDForUpdate(ctx, id)
		if err != nil {
			return err
		}

		trx, err := s.getActiveTransactionByCheckoutSessionIDForUpdate(ctx, session.ID)
		if err != nil {
			return err
		}

		if trx.Status == entity.TransactionStatusCompleted {
			return nil
		}
		if trx.Status != entity.TransactionStatusPending {
			return entity.ErrInvalidTransaction
		}

		err = s.repo.UpdateTransactionToCompletedByCheckoutSessionID(ctx, session.ID)
		if err != nil {
			slog.ErrorContext(ctx, "[StripeEventHandler-updatePendingTransactionToCompleted] fail update transaction", "error", err)
			return entity.ErrInternal
		}

		_, err = s.repo.AddWalletBalance(ctx, wallet.ID, trx.Amount)
		if err != nil {
			slog.ErrorContext(ctx, "[StripeEventHandler-updatePendingTransactionToCompleted] fail add balance to wallet", "error", err)
			return entity.ErrInternal
		}

		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "[StripeEventHandler-updatePendingTransactionToCompleted] fail tx manager", "error", err)
		return entity.ErrInternal
	}
	return nil
}
