package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
)

const (
	topupQuantity = 1
	topupPurpose  = "wallet topup"
)

// TopupWallet defines interface to topup wallet.
type TopupWallet interface {
	// Topup topups a wallet's balance.
	// It needs idempotency key.
	Topup(ctx context.Context, input *entity.TopupWalletInput) (*entity.TopupWalletOutput, error)
}

// TopupWalletRepository defines the interface to update wallet in repository.
type TopupWalletRepository interface {
	// GetActivePendingTransactionByIdempotencyKey gets a pending transaction by idempotency key.
	GetActivePendingTransactionByIdempotencyKey(ctx context.Context, key uuid.UUID) (*entity.Transaction, error)
	// GetActiveWalletByIDAndUserID gets a user's wallet.
	GetActiveWalletByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Wallet, error)
	// GetActiveCustomerByUserID gets a customer.
	GetActiveCustomerByUserID(ctx context.Context, userID uuid.UUID) (*entity.Customer, error)
	// InsertTransaction inserts a new transaction.
	InsertTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
}

// PaymentClient defines interface for payment.
type PaymentClient interface {
	// GetCheckoutSession gets a checkout session by id.
	GetCheckoutSession(ctx context.Context, id string) (*entity.CheckoutSession, error)
	// CreateCheckoutSession creates a checkout session.
	CreateCheckoutSession(ctx context.Context, input *entity.CheckoutInput) (*entity.CheckoutSession, error)
}

// WalletTopup is responsible for topup a wallet.
type WalletTopup struct {
	walletRepo    TopupWalletRepository
	paymentClient PaymentClient
	successURL    string
}

// NewWalletTopup creates an instance of WalletTopup.
func NewWalletTopup(w TopupWalletRepository, p PaymentClient, s string) *WalletTopup {
	return &WalletTopup{walletRepo: w, paymentClient: p, successURL: s}
}

// Topup tops up a wallet.
func (wt *WalletTopup) Topup(ctx context.Context, input *entity.TopupWalletInput) (*entity.TopupWalletOutput, error) {
	if err := validateTopupWalletInput(input); err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] topup input is invalid", "error", err)
		return nil, err
	}

	trx, err := wt.walletRepo.GetActivePendingTransactionByIdempotencyKey(ctx, input.IdempotencyKey)
	if err != nil && err != entity.ErrNilTransaction {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail get transaction", "error", err)
		return nil, entity.ErrInternal
	}
	// there is pending transaction with inputted idempotency key.
	// just return the payment checkout session.
	if trx != nil && trx.CheckoutSessionID != nil {
		var session *entity.CheckoutSession
		session, err = wt.paymentClient.GetCheckoutSession(ctx, *trx.CheckoutSessionID)
		if err != nil {
			slog.ErrorContext(ctx, "[WalletTopup-Topup] fail get checkout session", "error", err)
			return nil, entity.ErrInternal
		}
		return &entity.TopupWalletOutput{CheckoutSessionURL: session.URL}, nil
	}

	// all this flow below is for non-existent transaction
	wallet, err := wt.walletRepo.GetActiveWalletByIDAndUserID(ctx, input.WalletID, input.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail get wallet", "error", err)
		return nil, err
	}
	customer, err := wt.walletRepo.GetActiveCustomerByUserID(ctx, input.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail get customer", "error", err)
		return nil, err
	}

	checkoutInput := &entity.CheckoutInput{
		Amount:           input.Amount,
		Currency:         wallet.Currency,
		Quantity:         topupQuantity,
		StripeCustomerID: customer.StripeCustomerID,
		SuccessURL:       wt.successURL,
		Purpose:          topupPurpose,
		WalletID:         input.WalletID,
	}
	session, err := wt.paymentClient.CreateCheckoutSession(ctx, checkoutInput)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail create checkout session", "error", err)
		return nil, err
	}

	trx = createPendingTransaction(input, wallet.Currency, session.ID)
	_, err = wt.walletRepo.InsertTransaction(ctx, trx)
	if err != nil {
		slog.ErrorContext(ctx, "[WalletTopup-Topup] fail insert transaction", "error", err)
		return nil, err
	}
	return &entity.TopupWalletOutput{CheckoutSessionURL: session.URL}, nil
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

func createPendingTransaction(input *entity.TopupWalletInput, currency string, sessionID string) *entity.Transaction {
	trx := &entity.Transaction{
		ID:                uuid.Must(uuid.NewV7()),
		UserID:            input.UserID,
		Type:              entity.TransactionTypeTopup,
		Status:            entity.TransactionStatusPending,
		IdempotencyKey:    input.IdempotencyKey,
		Amount:            input.Amount,
		Currency:          currency,
		CheckoutSessionID: &sessionID,
	}
	setTransactionAuditableProperties(trx)
	return trx
}

func setTransactionAuditableProperties(transaction *entity.Transaction) {
	now := time.Now().UTC()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now
	transaction.CreatedBy = transaction.UserID
	transaction.UpdatedBy = transaction.UserID
}
